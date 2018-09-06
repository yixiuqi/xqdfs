package volume

import (
	"os"
	"fmt"
	"sync"
	"time"
	"sync/atomic"
	"path/filepath"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/stat"
	"xqdfs/storage/conf"
	"xqdfs/storage/block"
	"xqdfs/storage/index"
	"xqdfs/storage/needle"
	myos "xqdfs/storage/os"

	"github.com/BurntSushi/toml"
)

type Volume struct {
	lock sync.RWMutex
	// meta
	Id      int32             	`json:"id"`
	Stats   *stat.Stats       	`json:"stats"`
	Block   *block.SuperBlock 	`json:"block"`
	Indexer *index.Indexer    	`json:"index"`
	// data
	LastKey	int64				`json:"lastKey"`
	needles map[int64]int64
	conf    *conf.Config
	// compact
	Compact       bool   		`json:"compact"`
	CompactOffset uint32 		`json:"compactOffset"`
	CompactTime   int64  		`json:"compactTime"`
	// status
	closed bool
}

// NewVolume new a volume and init it.
func NewVolume(id int32, bfile, ifile string, c *conf.Config) (v *Volume, err error) {
	v = &Volume{
		Id:id,
		Stats:&stat.Stats{},
		needles:make(map[int64]int64),
		conf:c,
		Compact:false,
		CompactOffset:0,
		CompactTime:0,
		closed:false,
	}
	if v.Block, err = block.NewSuperBlock(bfile, c); err != nil {
		return nil, err
	}
	if v.Indexer, err = index.NewIndexer(ifile, c); err != nil {
		v.Close()
		return nil, err
	}
	if err = v.init(); err != nil {
		v.Close()
		return nil, err
	}
	return
}

func (v *Volume) loadStats(vid int32,ifile string) *stat.Stats {
	if vid==-1{
		return &stat.Stats{}
	}
	path:=fmt.Sprintf("%s/%d.toml",filepath.Dir(ifile),vid)
	log.Debugf("storeStats [%s] vid[%d]",path,vid)

	stat:=&stat.Stats{}
	if myos.Exist(path) == false {
		f, err:= os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
			return nil
		}
		defer f.Close()
		enc:=toml.NewEncoder(f)
		err=enc.Encode(stat)
		if err!=nil{
			log.Error(err)
			return nil
		}else{
			return stat
		}
	}else{
		f, err:= os.OpenFile(path, os.O_RDONLY, 0664)
		if err != nil {
			log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
			return nil
		}
		defer f.Close()

		_,err=toml.DecodeReader(f,stat)
		if err!=nil{
			log.Error(err)
			return nil
		}else{
			return stat
		}
	}
}

func (v *Volume) StoreStats() error {
	path:=fmt.Sprintf("%s/%d.toml",filepath.Dir(v.Indexer.File),v.Id)
	f, err:= os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
		return err
	}
	defer f.Close()
	if v.Stats==nil{
		return nil
	}

	enc:=toml.NewEncoder(f)
	err=enc.Encode(v.Stats)
	if err!=nil{
		log.Error(err)
		return err
	}else{
		return nil
	}
}

func (v *Volume) ImageCount() uint64 {
	return uint64(len(v.needles))
}

// Meta get index meta data.
func (v *Volume) Meta() []byte {
	return []byte(fmt.Sprintf("%s,%s,%d", v.Block.File, v.Indexer.File, v.Id))
}

// IsClosed reports whether the volume is closed.
func (v *Volume) IsClosed() bool {
	return v.closed
}

// init recovery super block from index or super block.
func (v *Volume) init() (err error) {
	var (
		size       int64
		offset     uint32
		lastOffset uint32
	)

	v.Stats = v.loadStats(v.Id,v.Indexer.File)

	// recovery from index
	if err = v.Indexer.Recovery(func(ix *index.Index) error {
		// must no less than last offset
		if ix.Offset < lastOffset {
			log.Errorf("recovery index: %s lastoffset: %d error(%s)", ix, lastOffset, "index offset error")
			return errors.ErrIndexOffset
		}
		// WARN if index's offset more than the block, discard it.
		if size = int64(ix.Size) + needle.BlockOffset(ix.Offset); size > v.Block.Size {
			log.Error("recovery index: %s EOF", ix)
			return errors.ErrIndexEOF
		}
		v.LastKey=ix.Key
		v.needles[ix.Key] = needle.NewCache(ix.Offset, ix.Size)
		offset = ix.Offset + needle.NeedleOffset(int64(ix.Size))
		lastOffset = ix.Offset
		return nil
	}); err != nil && err != errors.ErrIndexEOF {
		return
	}
	// recovery from super block
	if err = v.Block.Recovery(offset, func(n *needle.Needle, so, eo uint32) (err1 error) {
		if n.Flag == needle.FlagOK {
			if err1 = v.Indexer.Write(n.Key, so, n.TotalSize); err1 != nil {
				return
			}
		} else {
			so = needle.CacheDelOffset
		}
		v.LastKey=n.Key
		v.needles[n.Key] = needle.NewCache(so, n.TotalSize)
		return
	}); err != nil {
		return
	}
	// flush index
	err = v.Indexer.Flush()
	return
}

func (v *Volume) Open() (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.closed {
		return
	}
	if err = v.Block.Open(); err != nil {
		v.Close()
		return
	}
	if err = v.Indexer.Open(); err != nil {
		v.Close()
		return
	}
	if err = v.init(); err != nil {
		v.Close()
		return
	}
	v.closed = false
	return
}

func (v *Volume) close() {
	if v.Block != nil {
		v.Block.Close()
	}
	if v.Indexer != nil {
		v.Indexer.Close()
	}
	v.closed = true
}

// Close close the volume.
func (v *Volume) Close() {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.close()
}

func (v *Volume) Clear() (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	if v.closed {
		return
	}
	if v.Compact {
		err = errors.ErrVolumeInCompact
		return
	}

	if v.Block != nil {
		v.Block.Destroy()
	}
	if v.Indexer != nil {
		v.Indexer.Destroy()
	}

	v.Stats = &stat.Stats{}
	v.StoreStats()
	v.needles = make(map[int64]int64)
	v.LastKey=0
	v.Compact = false
	v.CompactOffset = 0
	v.CompactTime = 0
	v.closed = false
	if v.Block, err = block.NewSuperBlock(v.Block.File, v.conf); err != nil {
		return err
	}
	if v.Indexer, err = index.NewIndexer(v.Indexer.File, v.conf); err != nil {
		v.close()
		return err
	}
	if err = v.init(); err != nil {
		v.close()
		return err
	}
	return
}

// Destroy remove block and index file, must called after Close().
func (v *Volume) Destroy() {
	v.lock.Lock()
	defer v.lock.Unlock()
	if !v.closed {
		v.close()
	}
	if v.Block != nil {
		v.Block.Destroy()
	}
	if v.Indexer != nil {
		v.Indexer.Destroy()
	}
}

func (v *Volume) read(n *needle.Needle) (err error) {
	var (
		key  = n.Key
		size = n.TotalSize
	)
	// pread syscall is atomic, no lock
	if err = v.Block.ReadAt(n); err != nil {
		return
	}
	if n.Key != key {
		return errors.ErrNeedleKey
	}
	if n.TotalSize != size {
		return errors.ErrNeedleSize
	}

	// needles map may be out-dated, recheck
	if n.Flag == needle.FlagDel {
		v.lock.Lock()
		v.needles[key] = needle.NewCache(needle.CacheDelOffset, size)
		v.lock.Unlock()
		err = errors.ErrNeedleDeleted
	}else {
		atomic.AddUint64(&v.Stats.TotalReadProcessed, 1)
		atomic.AddUint64(&v.Stats.TotalReadBytes, uint64(size))
	}
	return
}

func (v *Volume) GetHeader(key int64) (n *needle.Needle, err error) {
	var (
		ok bool
		nc int64
	)
	v.lock.RLock()
	if nc, ok = v.needles[key]; !ok {
		err = errors.ErrNeedleNotExist
	}
	v.lock.RUnlock()
	if err == nil {
		offset, _:= needle.Cache(nc)
		if offset != needle.CacheDelOffset {
			n,err= v.Block.GetHeader(offset)
		} else {
			err = errors.ErrNeedleDeleted
		}
	}
	return
}

// Read get a needle by key and cookie and write to wr.
func (v *Volume) Read(key int64, cookie int32) (n *needle.Needle, err error) {
	var (
		ok bool
		nc int64
	)
	v.lock.RLock()
	if nc, ok = v.needles[key]; !ok {
		err = errors.ErrNeedleNotExist
	}
	v.lock.RUnlock()
	if err == nil {
		if n = needle.NewReader(key, nc); n.Offset != needle.CacheDelOffset {
			if err = v.read(n); err == nil {
				if n.Cookie != cookie {
					err = errors.ErrNeedleCookie
				}
			}
		} else {
			err = errors.ErrNeedleDeleted
		}
		if err != nil {
			n.Close()
			n = nil
		}
	}
	return
}

func (v *Volume) Write(n *needle.Needle) (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	var (
		ok     bool
	)
	_, ok = v.needles[n.Key]
	if ok {
		err=errors.ErrNeedleExist
		v.lock.Unlock()
		return
	}

	n.Offset = v.Block.Offset
	if err = v.Block.Write(n); err == nil {
		if err = v.Indexer.Write(n.Key, n.Offset, n.TotalSize); err == nil {
			v.LastKey=n.Key
			v.needles[n.Key] = needle.NewCache(n.Offset, n.TotalSize)
		}
	}

	if err == nil {
		atomic.AddUint64(&v.Stats.TotalWriteProcessed, 1)
		atomic.AddUint64(&v.Stats.TotalWriteBytes, uint64(n.TotalSize))
	}
	return
}

func (v *Volume) Rewrite(n *needle.Needle) (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	var (
		ok     bool
		nc     int64
		offset uint32
	)
	n.Offset = v.Block.Offset
	if err = v.Block.Write(n); err == nil {
		if err = v.Indexer.Write(n.Key, n.Offset, n.TotalSize); err == nil {
			v.LastKey=n.Key
			nc, ok = v.needles[n.Key]
			v.needles[n.Key] = needle.NewCache(n.Offset, n.TotalSize)
		}
	}

	if err == nil {
		if ok {
			log.Debug("needle is rewrite,key:",n.Key)
			offset, _ = needle.Cache(nc)
			v.del(offset)
		}
		atomic.AddUint64(&v.Stats.TotalWriteProcessed, 1)
		atomic.AddUint64(&v.Stats.TotalWriteBytes, uint64(n.TotalSize))
	}
	return
}

// del signal the godel goroutine aync merge all offsets and del.
func (v *Volume) del(offset uint32) (err error) {
	if offset == needle.CacheDelOffset {
		return
	}
	if err = v.Block.Delete(offset); err != nil {
		log.Error("volume delete error")
		return
	}
	atomic.AddUint64(&v.Stats.TotalDelProcessed, 1)
	atomic.AddUint64(&v.Stats.TotalWriteBytes, 1)
	return
}

func (v *Volume) Delete(key int64) (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	var (
		ok     bool
		nc     int64
		size   int32
		offset uint32
	)
	if nc, ok = v.needles[key]; ok {
		if offset, size = needle.Cache(nc); offset != needle.CacheDelOffset {
			v.needles[key] = needle.NewCache(needle.CacheDelOffset, size)
		} else {
			err = errors.ErrNeedleDeleted
		}
	} else {
		err = errors.ErrNeedleNotExist
	}

	if err == nil {
		err = v.del(offset)
	}
	return
}

// compact compact v to new v.
func (v *Volume) compact(nv *Volume) (err error) {
	err = v.Block.Compact(v.CompactOffset, func(n *needle.Needle, so, eo uint32) (err1 error) {
		if n.Flag != needle.FlagDel {
			if err1 = nv.Write(n); err1 != nil {
				log.Error(err1)
				return
			}
		}
		v.CompactOffset = eo
		return
	})
	return
}

// Compact copy the super block to another space, and drop the "delete"
// needle, so this can reduce disk space cost.
func (v *Volume) StartCompact(nv *Volume) (err error) {
	v.lock.Lock()
	if v.Compact {
		err = errors.ErrVolumeInCompact
	} else {
		v.Compact = true
	}
	v.lock.Unlock()
	if err != nil {
		return
	}
	v.CompactTime = time.Now().UnixNano()
	if err = v.compact(nv); err != nil {
		return
	}
	return
}

func (v *Volume) StopCompact(nv *Volume) (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	if nv != nil {
		if err = v.compact(nv); err != nil {
			log.Error(err)
			goto free
		}

		v.Block, nv.Block = nv.Block, v.Block
		v.Indexer, nv.Indexer = nv.Indexer, v.Indexer
		v.needles, nv.needles = nv.needles, v.needles

		atomic.StoreUint64(&v.Stats.TotalDelProcessed,0)
		v.StoreStats()
	}
free:
	v.Compact = false
	v.CompactOffset = 0
	v.CompactTime = 0
	return
}

func (v *Volume) IsCompact() bool {
	return v.Compact
}
