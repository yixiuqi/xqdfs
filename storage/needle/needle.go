package needle

import (
	"hash/crc32"
	"syscall"
	"sync"
	"bytes"
	"bufio"
	"fmt"
	"io"

	"xqdfs/errors"
	"xqdfs/utils/encoding/binary"
)

// Needle stored int super block, aligned to 8bytes.
//
// needle file format:
//  --------------
// | super  block |
//  --------------
// |    needle    |		    ----------------
// |    needle    |        |  magic (int32) |
// |    needle    | ---->  |  cookie (int32)|
// |    needle    |        |  key (int64)   |
// |    needle    |        |  flag (byte)   |
// |    needle    |        |  size (int32)  |
// |    needle    |        |  data (bytes)  |
// |    needle    |        |  magic (int32) |
// |    needle    |        | checksum(int32)|
// |    needle    |        | padding (bytes)|
// |    ......    |         ----------------
// |    ......    |             int bigendian
//
// field     | explanation
// ---------------------------------------------------------
// magic     | header magic number used for checksum
// cookie    | random number to mitigate brute force lookups
// key       | 64bit photo id
// flag      | signifies deleted status
// size      | data size
// data      | the actual photo data
// magic     | footer magic number used for checksum
// checksum  | used to check integrity
// padding   | total needle size is aligned to 8 bytes

const (
// size
// footer
	MagicSize  	= 	4
	CookieSize 	= 	4
	KeySize    	= 	8
	FlagSize   	= 	1
	SizeSize   	= 	4
// data
// footer
// magic
	ChecksumSize = 4
// padding

// offset
// header
	MagicOffset  	= 	0
	CookieOffset 	= 	MagicOffset + MagicSize
	KeyOffset		= 	CookieOffset + CookieSize
	FlagOffset   	= 	KeyOffset + KeySize
	SizeOffset   	= 	FlagOffset + FlagSize
	DataOffset   	= 	SizeOffset + SizeSize

// footer
// MagicOffset  	= 0
	ChecksumOffset 	= MagicOffset + MagicSize
	PaddingOffset  	= ChecksumOffset + ChecksumSize

// header is constant = 21
	HeaderSize = MagicSize + CookieSize + KeySize + FlagSize + SizeSize
// footer is constant = 8 (no padding)
	FooterSize = MagicSize + ChecksumSize

// WARN our offset is aligned with padding size(8)
// so a uint32 can store 4GB * 8 offset
// if you want a block more larger, modify this constant, but must bigger
// than 8
	PaddingSize   = 8
	PaddingAlign = PaddingSize - 1
	PaddingByte  = byte(0)

// flags
	FlagOK  = byte(0)
	FlagDel = byte(1)

// display
	DisplayData = 16
)

var (
	Padding = [][]byte{nil}
	// crc32 checksum table, goroutine safe
	Crc32Table = crc32.MakeTable(crc32.Koopman)
	// magic number
	HeaderMagic = []byte{12, 34, 56, 78}
	FooterMagic = []byte{87, 65, 43, 21}
	// flag
	FlagDelBytes = []byte{FlagDel}

	PageSize = syscall.Getpagesize()
	BufPool  = sync.Pool{
		New: func() interface{} {
			return make([]byte, PageSize) // 4kb
		},
	}
)

func init() {
	var i int
	for i = 1; i < PaddingSize; i++ {
		Padding = append(Padding, bytes.Repeat([]byte{PaddingByte}, i))
	}
	return
}

type Needle struct {
	HeaderMagic []byte	`json:"-"`
	Cookie      int32	`json:"cookie"`
	Key         int64	`json:"key"`
	Flag        byte	`json:"flag"`
	Size        int32	`json:"size"`
	Data        []byte	`json:"-"`
	FooterMagic []byte	`json:"-"`
	Checksum    uint32	`json:"-"`
	Padding     []byte	`json:"-"`

	PaddingSize int32	`json:"paddingSize"`
	TotalSize   int32	`json:"totalSize"`
	FooterSize  int32	`json:"footerSize"`
	IncrOffset uint32	`json:"incrOffset"`		//NeedleOffset(int64(n.TotalSize))
	Offset     uint32	`json:"offset"`
	buffer     []byte
}

func NewWriter(key int64, cookie, size int32) *Needle {
	var n = new(Needle)
	n.Key = key
	n.Cookie = cookie
	n.Size = size
	n.init()
	n.newBuffer()
	return n
}

func (n *Needle) InitWriter(key int64, cookie, size int32) {
	n.Key = key
	n.Cookie = cookie
	n.Size = size
	n.init()
	n.newBuffer()
}

func NewReader(key, nc int64) *Needle {
	var n = new(Needle)
	n.Key = key
	n.Offset, n.TotalSize = Cache(nc)
	n.newBuffer()
	return n
}

func (n *Needle) Close() {
	n.freeBuffer()
}

func (n *Needle) newBuffer() {
	if n.TotalSize <= int32(PageSize) {
		n.buffer = BufPool.Get().([]byte)
	} else {
		n.buffer = make([]byte, n.TotalSize)
	}
}

func (n *Needle) freeBuffer() {
	if n.buffer != nil && len(n.buffer) <= PageSize {
		BufPool.Put(n.buffer)
	}
}

func (n *Needle) Buffer() []byte {
	return n.buffer[:n.TotalSize]
}

func (n *Needle) calcSize() {
	n.TotalSize = int32(HeaderSize + n.Size + FooterSize)
	n.PaddingSize = align(n.TotalSize) - n.TotalSize
	n.TotalSize += n.PaddingSize
	n.FooterSize = FooterSize + n.PaddingSize
	n.IncrOffset = NeedleOffset(int64(n.TotalSize))
}

func (n *Needle) ParseHeader(buf []byte) (err error) {
	if len(buf) != HeaderSize {
		return errors.ErrNeedleHeaderSize
	}
	// magic
	n.HeaderMagic = buf[MagicOffset:CookieOffset]
	if !bytes.Equal(n.HeaderMagic, HeaderMagic) {
		return errors.ErrNeedleHeaderMagic
	}
	// cookie
	n.Cookie = binary.BigEndian.Int32(buf[CookieOffset:KeyOffset])
	// key
	n.Key = binary.BigEndian.Int64(buf[KeyOffset:FlagOffset])
	// flag
	n.Flag = buf[FlagOffset]
	if n.Flag != FlagOK && n.Flag != FlagDel {
		return errors.ErrNeedleFlag
	}
	// size
	n.Size = binary.BigEndian.Int32(buf[SizeOffset:DataOffset])
	if n.Size < 0 {
		return errors.ErrNeedleSize
	}
	n.calcSize()
	return
}

func (n *Needle) parseData(buf []byte) (err error) {
	if len(buf) != int(n.Size) {
		return errors.ErrNeedleDataSize
	}
	// data
	n.Data = buf
	// checksum
	n.Checksum = crc32.Update(0, Crc32Table, n.Data)
	return
}

func (n *Needle) parseFooter(buf []byte) (err error) {
	if len(buf) != int(n.FooterSize) {
		return errors.ErrNeedleFooterSize
	}
	// magic
	n.FooterMagic = buf[MagicOffset:ChecksumOffset]
	if !bytes.Equal(n.FooterMagic, FooterMagic) {
		return errors.ErrNeedleFooterMagic
	}
	if n.Checksum != binary.BigEndian.Uint32(buf[ChecksumOffset:PaddingOffset]) {
		return errors.ErrNeedleChecksum
	}
	// padding
	n.Padding = buf[PaddingOffset : PaddingOffset+n.PaddingSize]
	if !bytes.Equal(n.Padding, Padding[n.PaddingSize]) {
		return errors.ErrNeedlePadding
	}
	return
}

func (n *Needle) writeHeader(buf []byte) (err error) {
	if len(buf) != int(HeaderSize) {
		return errors.ErrNeedleHeaderSize
	}
	// magic
	copy(buf[MagicOffset:CookieOffset], n.HeaderMagic)
	// cookie
	binary.BigEndian.PutInt32(buf[CookieOffset:KeyOffset], n.Cookie)
	// key
	binary.BigEndian.PutInt64(buf[KeyOffset:FlagOffset], n.Key)
	// flag
	buf[FlagOffset] = n.Flag
	// size
	binary.BigEndian.PutInt32(buf[SizeOffset:DataOffset], n.Size)
	return
}

func (n *Needle) writeFooter(buf []byte) (err error) {
	if len(buf) != int(n.FooterSize) {
		return errors.ErrNeedleFooterSize
	}
	// magic
	copy(buf[MagicOffset:ChecksumOffset], n.FooterMagic)
	// checksum
	binary.BigEndian.PutUint32(buf[ChecksumOffset:PaddingOffset], n.Checksum)
	// padding
	copy(buf[PaddingOffset:PaddingOffset+n.PaddingSize], n.Padding)
	return
}

func (n *Needle) ParseFrom(rd *bufio.Reader) (err error) {
	var (
		dataOffset   int32
		footerOffset int32
		endOffset    int32
		data         []byte
	)
	// header
	if data, err = rd.Peek(HeaderSize); err != nil {
		return
	}
	if err = n.ParseHeader(data); err != nil {
		return
	}
	dataOffset = HeaderSize
	footerOffset = dataOffset + n.Size
	endOffset = footerOffset + n.FooterSize
	// no discard, get all needle buffer
	if data, err = rd.Peek(int(n.TotalSize)); err != nil {
		return
	}
	if err = n.parseData(data[dataOffset:footerOffset]); err != nil {
		return
	}
	// footer
	if err = n.parseFooter(data[footerOffset:endOffset]); err != nil {
		return
	}
	n.buffer = data
	_, err = rd.Discard(int(n.TotalSize))
	return
}

func (n *Needle) Parse() (err error) {
	var dataOffset int32
	if err = n.ParseHeader(n.buffer[:HeaderSize]); err == nil {
		dataOffset = HeaderSize + n.Size
		if err = n.parseData(n.buffer[HeaderSize:dataOffset]); err == nil {
			err = n.parseFooter(n.buffer[dataOffset:n.TotalSize])
		}
	}
	return
}

func (n *Needle) ReadFrom(rd io.Reader) (err error) {
	var (
		dataOffset int32
		data       []byte
	)
	dataOffset = HeaderSize + n.Size
	data = n.buffer[HeaderSize:dataOffset]
	if err = n.writeHeader(n.buffer[:HeaderSize]); err == nil {
		if _, err = rd.Read(data); err == nil {
			n.Data = data
			n.Checksum = crc32.Update(0, Crc32Table, data)
			err = n.writeFooter(n.buffer[dataOffset:n.TotalSize])
		}
	}
	return
}

func (n *Needle) String() string {
	var dn = DisplayData
	if len(n.Data) < dn {
		dn = len(n.Data)
	}
	return fmt.Sprintf(`
-----------------------------
TotalSize[%d]
HeaderSize[%d],HeaderMagic[%v],Cookie[%d],Key[%d],Flag[%d],Size[%d]
Data[%v...]
FooterSize[%d],FooterMagic[%v],Checksum[%d],Padding[%v]
-----------------------------`,
	n.TotalSize, HeaderSize, n.HeaderMagic, n.Cookie, n.Key, n.Flag, n.Size,n.Data[:dn], n.FooterSize, n.FooterMagic, n.Checksum, n.Padding)
}

func (n *Needle) init() {
	n.calcSize()
	n.Flag = FlagOK
	n.HeaderMagic = HeaderMagic
	n.FooterMagic = FooterMagic
	n.Padding = Padding[n.PaddingSize]
	return
}



