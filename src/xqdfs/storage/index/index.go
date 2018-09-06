package index

import (
	"os"
	"fmt"
	"bufio"
	"io"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/storage/conf"
	"xqdfs/storage/needle"
	"xqdfs/utils/encoding/binary"
	myos "xqdfs/storage/os"
)

// Index for fast recovery super block needle cache in memory, index is async
// append the needle meta data.
//
// index file format:
//  ---------------
// | super   block |
//  ---------------
// |     needle    |		   ----------------
// |     needle    |          |  key (int64)   |
// |     needle    | ---->    |  offset (uint) |
// |     needle    |          |  size (int32)  |
// |     ......    |           ----------------
// |     ......    |             int bigendian
//
// field     | explanation
// --------------------------------------------------
// key       | needle key (photo id)
// offset    | needle offset in super block (aligned)
// size      | needle data size

const (
	// signal command
	Finish = 0
	Ready  = 1
	// index size
	KeySize    = 8
	OffsetSize = 4
	SizeSize   = 4
	// index size = 16
	IndexSize = KeySize + OffsetSize + SizeSize
	// index offset
	KeyOffset    = 0
	OffsetOffset = KeyOffset + KeySize
	SizeOffset   = OffsetOffset + OffsetSize
	// 100mb
	FallocSize = 100 * 1024 * 1024
)

type Indexer struct {
	f      *os.File
	// buffer
	buf []byte
	bn  int
	//
	File    string `json:"file"`
	LastErr error  `json:"-"`
	Offset  int64  `json:"offset"`
	conf    *conf.Config
	// status
	syncOffset int64
	closed     bool
	write      int
}

type Index struct {
	Key    int64
	Offset uint32	//needle offset
	Size   int32
}

func NewIndexer(file string, conf *conf.Config) (i *Indexer, err error) {
	var stat os.FileInfo
	i = &Indexer{}
	i.File = file
	i.closed = false
	i.syncOffset = 0
	i.conf = conf
	i.bn = 0
	if conf.Index.BufferSize < IndexSize {
		i.buf = make([]byte, IndexSize)
	} else {
		i.buf = make([]byte, conf.Index.BufferSize)
	}
	if i.f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|myos.O_NOATIME, 0664); err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", file, err)
		return nil, err
	}
	if stat, err = i.f.Stat(); err != nil {
		log.Errorf("index: %s Stat() error(%v)", i.File, err)
		return nil, err
	}
	if stat.Size() == 0 {
		if err = myos.Fallocate(i.f.Fd(), myos.FALLOC_FL_KEEP_SIZE, 0, FallocSize); err != nil {
			log.Errorf("index: %s fallocate() error(err)", i.File, err)
			i.Close()
			return nil, err
		}
	}
	return
}

func (i *Index) parse(buf []byte) (err error) {
	i.Key = binary.BigEndian.Int64(buf)
	i.Offset = binary.BigEndian.Uint32(buf[OffsetOffset:])
	i.Size = binary.BigEndian.Int32(buf[SizeOffset:])
	if i.Size < 0 {
		return errors.ErrIndexSize
	}
	return
}

func (i *Index) String() string {
	return fmt.Sprintf(`Key[%d],Offset[%d],Size[%d]`, i.Key, i.Offset, i.Size)
}

// Write append index needle to disk.
// WARN can't concurrency with merge and write.
// ONLY used in super block recovery!!!!!!!!!!!
func (i *Indexer) Write(key int64, offset uint32, size int32) (err error) {
	if i.LastErr != nil {
		return i.LastErr
	}
	if i.bn+IndexSize >= i.conf.Index.BufferSize {
		// buffer full
		if err = i.flush(true); err != nil {
			return
		}
	}
	binary.BigEndian.PutInt64(i.buf[i.bn:], key)
	i.bn += KeySize
	binary.BigEndian.PutUint32(i.buf[i.bn:], offset)
	i.bn += OffsetSize
	binary.BigEndian.PutInt32(i.buf[i.bn:], size)
	i.bn += SizeSize
	err = i.flush(false)
	return
}

// flush the in-memory data flush to disk.
func (i *Indexer) flush(force bool) (err error) {
	var (
		fd     uintptr
		offset int64
		size   int64
	)
	if i.write++; !force && i.write < i.conf.Index.SyncWrite {
		return
	}
	if _, err = i.f.Write(i.buf[:i.bn]); err != nil {
		i.LastErr = err
		log.Errorf("index: %s Write() error(%v)", i.File, err)
		return
	}
	i.Offset += int64(i.bn)
	i.bn = 0
	i.write = 0
	offset = i.syncOffset
	size = i.Offset - i.syncOffset
	fd = i.f.Fd()
	if i.conf.Index.Syncfilerange {
		if err = myos.Syncfilerange(fd, offset, size, myos.SYNC_FILE_RANGE_WRITE); err != nil {
			i.LastErr = err
			log.Errorf("index: %s Syncfilerange() error(%v)", i.File, err)
			return
		}
	} else {
		if err = myos.Fdatasync(fd); err != nil {
			i.LastErr = err
			log.Errorf("index: %s Fdatasync() error(%v)", i.File, err)
			return
		}
	}
	if err = myos.Fadvise(fd, offset, size, myos.POSIX_FADV_DONTNEED); err == nil {
		i.syncOffset = i.Offset
	} else {
		log.Errorf("index: %s Fadvise() error(%v)", i.File, err)
		i.LastErr = err
	}
	return
}

// Flush flush writer buffer.
func (i *Indexer) Flush() (err error) {
	if i.LastErr != nil {
		return i.LastErr
	}
	err = i.flush(true)
	return
}

// Scan scan a indexer file.
func (i *Indexer) Scan(r *os.File, fn func(*Index) error) (err error) {
	var (
		data []byte
		fi   os.FileInfo
		fd   = r.Fd()
		ix   = &Index{}
		rd   = bufio.NewReaderSize(r, i.conf.Index.BufferSize)
	)
	log.Debugf("scan index: %s", i.File)
	// advise sequential read
	if fi, err = r.Stat(); err != nil {
		log.Errorf("index: %s Stat() error(%v)", i.File)
		return
	}
	if err = myos.Fadvise(fd, 0, fi.Size(), myos.POSIX_FADV_SEQUENTIAL); err != nil {
		log.Errorf("index: %s Fadvise() error(%v)", i.File)
		return
	}
	if _, err = r.Seek(0, os.SEEK_SET); err != nil {
		log.Errorf("index: %s Seek() error(%v)", i.File, err)
		return
	}
	for {
		if data, err = rd.Peek(IndexSize); err != nil {
			break
		}
		if err = ix.parse(data); err != nil {
			break
		}
		if ix.Size > int32(needle.Size(i.conf.NeedleMaxSize)) {
			log.Errorf("scan index: %s error(%s)", ix, "index size error")
			err = errors.ErrIndexSize
			break
		}
		if _, err = rd.Discard(IndexSize); err != nil {
			break
		}

		if err = fn(ix); err != nil {
			break
		}
	}
	if err == io.EOF {
		// advise no need page cache
		if err = myos.Fadvise(fd, 0, fi.Size(), myos.POSIX_FADV_DONTNEED); err == nil {
			err = nil
			return
		} else {
			log.Errorf("index: %s Fadvise() error(%v)", i.File)
		}
	}
	return
}

// Recovery recovery needle cache meta data in memory, index file  will stop
// at the right parse data offset.
func (i *Indexer) Recovery(fn func(*Index) error) (err error) {
	if i.Scan(i.f, func(ix *Index) (err1 error) {
		if err1 = fn(ix); err1 == nil {
			i.Offset += int64(IndexSize)
		}
		return
	}); err != nil {
		return
	}
	// reset b.w offset, discard left space which can't parse to a needle
	if _, err = i.f.Seek(i.Offset, os.SEEK_SET); err != nil {
		log.Errorf("index: %s Seek() error(%v)", i.File, err)
	}
	return
}

// Open open the closed indexer, must called after NewIndexer.
func (i *Indexer) Open() (err error) {
	if !i.closed {
		return
	}
	if i.f, err = os.OpenFile(i.File, os.O_RDWR|myos.O_NOATIME, 0664); err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", i.File, err)
		return
	}
	// reset buf
	i.bn = 0
	i.closed = false
	i.LastErr = nil
	return
}

// Close close the indexer file.
func (i *Indexer) Close() {
	var err error
	if i.f != nil {
		if err = i.flush(true); err != nil {
			log.Errorf("index: %s Flush() error(%v)", i.File, err)
		}
		if err = i.f.Sync(); err != nil {
			log.Errorf("index: %s Sync() error(%v)", i.File, err)
		}
		if err = i.f.Close(); err != nil {
			log.Errorf("index: %s Close() error(%v)", i.File, err)
		}
	}
	i.closed = true
	i.LastErr = errors.ErrIndexClosed
	return
}

// Destroy destroy the indexer.
func (i *Indexer) Destroy() {
	if !i.closed {
		i.Close()
	}
	os.Remove(i.File)
}