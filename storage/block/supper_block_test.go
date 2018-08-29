package block

import (
	"xqdfs/storage/conf"
	"xqdfs/storage/needle"
	"testing"
	"os"
	"fmt"
	"bytes"
	"math/rand"
)

var (
	testConf = &conf.Config{
		NeedleMaxSize: 4 * 1024 * 1024,
		Block: &conf.Block{
			BufferSize:    4 * 1024 * 1024,
			SyncWrite:     1,
			Syncfilerange: true,
		},
	}
)

func TestSuperBlock(t *testing.T) {
	var (
		b                  *SuperBlock
		n                  *needle.Needle
		//offset, v2, v3, v4 uint32
		err                error
		buf                = &bytes.Buffer{}
		//needles            = make(map[int64]int64)
		data               = []byte("test")
		file               = "./test.block"
	)

	os.Remove(file)
	defer os.Remove(file)
	{
		fmt.Println("--------------------------------------------------------------------------- test new block file")
		if b, err = NewSuperBlock(file, testConf); err != nil {
			t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
			t.FailNow()
		}
		fmt.Println(b.String())
		b.Close()
	}

	{
		fmt.Println("--------------------------------------------------------------------------- test parse block file")
		if b, err = NewSuperBlock(file, testConf); err != nil {
			t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
			t.FailNow()
		}
		fmt.Println(b.String())
		b.Close()
	}

	{
		fmt.Println("--------------------------------------------------------------------------- test open")
		if err = b.Open(); err != nil {
			t.Errorf("Open() error(%v)", err)
			t.FailNow()
		}
		fmt.Println(b.String())
		b.Close()
	}

	{
		fmt.Println("--------------------------------------------------------------------------- test write")
		if b, err = NewSuperBlock(file, testConf); err != nil {
			t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
			t.FailNow()
		}
		fmt.Println(b.String())
		defer b.Close()

		//
		dataLen:=len(data)
		if _, err = buf.Write(data); err != nil {
			t.Errorf("buf.Write() error(%v)", err)
			t.FailNow()
		}
		n = needle.NewWriter(1, 2, int32(dataLen))
		defer n.Close()
		n.ReadFrom(buf)
		fmt.Println(n.String())

		//
		if err = b.Write(n); err != nil {
			t.Errorf("b.Write() error(%v)", err)
			t.FailNow()
		}

		if err = compareTestOffset(b, n, needle.NeedleOffset(int64(HeaderSize))); err != nil {
			t.Errorf("compareTestOffset() error(%v)", err)
			t.FailNow()
		}

		fmt.Println(b.String())
	}

	{
		fmt.Println("--------------------------------------------------------------------------- test get")
		if b, err = NewSuperBlock(file, testConf); err != nil {
			t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
			t.FailNow()
		}
		fmt.Println(b.String())
		defer b.Close()

		n = needle.NewReader(1, needle.NewCache(1,40))
		defer n.Close()

		if err = b.ReadAt(n); err != nil {
			t.Errorf("b.ReadAt() error(%v)", err)
			t.FailNow()
		}
		fmt.Println(n.String())

		if err = compareTestNeedle(t, 1, 2, needle.FlagOK, n, data); err != nil {
			t.Errorf("compareTestNeedle() error(%v)", err)
			t.FailNow()
		}
	}

	{
		fmt.Println("--------------------------------------------------------------------------- test delete")
		if b, err = NewSuperBlock(file, testConf); err != nil {
			t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
			t.FailNow()
		}
		fmt.Println(b.String())
		defer b.Close()

		//
		dataLen:=len(data)
		if _, err = buf.Write(data); err != nil {
			t.Errorf("buf.Write() error(%v)", err)
			t.FailNow()
		}
		n = needle.NewWriter(1, 2, int32(dataLen))
		defer n.Close()
		n.ReadFrom(buf)
		fmt.Println(n.String())

		//
		if err = b.Write(n); err != nil {
			t.Errorf("b.Write() error(%v)", err)
			t.FailNow()
		}

		if err = b.Delete(1); err != nil {
			t.Errorf("Del() error(%v)", err)
			t.FailNow()
		}

		//read
		n.Offset=1
		if err = b.ReadAt(n); err != nil {
			t.Errorf("b.ReadAt() error(%v)", err)
			t.FailNow()
		}
		fmt.Println(n.String())
	}
}

func compareTestOffset(b *SuperBlock, n *needle.Needle, offset uint32) (err error) {
	var v int64
	if b.Offset != offset+needle.NeedleOffset(int64(n.TotalSize)) {
		err = fmt.Errorf("b.Offset: %d not match %d", b.Offset, offset)
		return
	}
	if v, err = b.w.Seek(0, os.SEEK_CUR); err != nil {
		err = fmt.Errorf("b.Seek() error(%v)", err)
		return
	} else {
		if v != needle.BlockOffset(b.Offset) {
			err = fmt.Errorf("offset: %d not match", v)
			return
		}
	}
	return
}

func compareTestNeedle(t *testing.T, key int64, cookie int32, flag byte, n *needle.Needle, data []byte) (err error) {
	if !bytes.Equal(n.Data, data) {
		err = fmt.Errorf("data: %s not match", n.Data)
		t.Error(err)
		return
	}
	if n.Cookie != cookie {
		err = fmt.Errorf("cookie: %d not match", n.Cookie)
		t.Error(err)
		return
	}
	if n.Key != key {
		err = fmt.Errorf("key: %d not match", n.Key)
		t.Error(err)
		return
	}
	if n.Flag != flag {
		err = fmt.Errorf("flag: %d not match", n.Flag)
		t.Error(err)
		return
	}
	if n.Size != int32(len(data)) {
		err = fmt.Errorf("size: %d not match", n.Size)
		t.Error(err)
		return
	}
	return
}

func TestSuperBlockWrite1(t *testing.T) {
	fmt.Println("TestSuperBlockWrite1")

	var (
		b                  *SuperBlock
		n                  *needle.Needle
		err                error
		file               = "./test.block"
	)

	os.Remove(file)
	if b, err = NewSuperBlock(file, testConf); err != nil {
		t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
		t.FailNow()
	}
	defer b.Close()

	data:=bytes.Repeat([]byte{1},128)
	buf:= &bytes.Buffer{}
	if _, err = buf.Write(data); err != nil {
		t.Errorf("buf.Write() error(%v)", err)
		t.FailNow()
	}

	n = needle.NewWriter(1, 345, int32(len(data)))
	defer n.Close()

	n.ReadFrom(buf)
	if err = b.Write(n); err != nil {
		t.Errorf("b.Write() error(%v)", err)
		t.FailNow()
	}
	fmt.Println(b)
}

func TestSuperBlockWrite(t *testing.T) {
	fmt.Println("TestSuperBlockWrite")

	var (
		b                  *SuperBlock
		n                  *needle.Needle
		err                error
		file               = "./test.block"
	)

	os.Remove(file)
	if b, err = NewSuperBlock(file, testConf); err != nil {
		t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
		t.FailNow()
	}
	defer b.Close()

	data:=bytes.Repeat([]byte{1},rand.Intn(1024*64)+1024*64)
	for i:=0;i<10000;i++{
		if (i%1000)==0{
			fmt.Println(i+1)
		}

		func(){
			buf:= &bytes.Buffer{}
			if _, err = buf.Write(data); err != nil {
				t.Errorf("buf.Write() error(%v)", err)
				t.FailNow()
			}

			n = needle.NewWriter(int64(i)+1, 345, int32(len(data)))
			defer n.Close()

			n.ReadFrom(buf)
			if err = b.Write(n); err != nil {
				t.Errorf("b.Write() error(%v)", err)
				t.FailNow()
			}
		}()
	}

	fmt.Println(b)
}

func TestSuperBlockRead(t *testing.T) {
	fmt.Println("TestSuperBlockRead")

	var (
		b                  *SuperBlock
		err                error
		file               = "./test.block"
	)

	if b, err = NewSuperBlock(file, testConf); err != nil {
		t.Errorf("NewSuperBlock(\"%s\") error(%v)", file, err)
		t.FailNow()
	}
	fmt.Println(b.String())
	defer b.Close()

	//
	count:=0
	allLen:=int64(0)
	if err = b.Recovery(0, func(rn *needle.Needle, so, eo uint32) (err1 error) {
		if rn.Flag != needle.FlagOK {
			so = needle.CacheDelOffset
		}
		count++
		allLen+=int64(rn.TotalSize)
		return
	}); err != nil {
		t.Errorf("Recovery() error(%v)", err)
		t.FailNow()
	}

	fmt.Println("count:",count,"allLen:",allLen)
}