package needle

import (
	"testing"
	"fmt"
	"bytes"
	"hash/crc32"
	"bufio"
)

func TestNeedlePageSize(t *testing.T) {
	fmt.Println(PageSize)
}

func TestNeedlePadding(t *testing.T) {
	fmt.Println(Padding)
}

func TestNeedleReadFrom(t *testing.T) {
	var (
		err       error
		n	      *Needle
		data1     = []byte("tes1")
		buf       = &bytes.Buffer{}
		checksum1 = crc32.Update(0, Crc32Table, data1)
	)

	if _, err = buf.Write(data1); err != nil {
		t.Error(err)
		t.FailNow()
	}

	n = NewWriter(1, 2, 4)
	defer n.Close()

	if err = n.ReadFrom(buf); err != nil {
		t.Error(err)
		t.FailNow()
	}

	compareNeedle(t, n, 1, 2, data1, FlagOK, checksum1)
}

func TestNeedleParseFrom(t *testing.T) {
	var (
		br        *bufio.Reader
		err       error
		writer	  *Needle
		reader	  *Needle
		data1     = []byte("tes1")
		bufW       = &bytes.Buffer{}
		bufR       = &bytes.Buffer{}
		checksum1 = crc32.Update(0, Crc32Table, data1)
	)

	writer = NewWriter(1, 2, 4)
	defer writer.Close()
	bufW.Write(data1)
	writer.ReadFrom(bufW)

	bufR.Write(writer.Buffer())

	br = bufio.NewReader(bufR)
	reader = new(Needle)
	if err = reader.ParseFrom(br); err != nil {
		t.Error(err)
		t.FailNow()
	}
	compareNeedle(t, reader, 1, 2, data1, FlagOK, checksum1)
}

func TestNeedleOffset(t *testing.T) {
	var (
		offset  int64
		noffset uint32
	)
	offset = 32
	if noffset = NeedleOffset(offset); noffset != uint32(offset/int64(PaddingSize)) {
		t.Errorf("noffset: %d not match", noffset)
		t.FailNow()
	}else{
		t.Log(noffset)
	}
	offset = 48
	if noffset = NeedleOffset(offset); noffset != uint32(offset/int64(PaddingSize)) {
		t.Errorf("noffset: %d not match", noffset)
		t.FailNow()
	}else{
		t.Log(noffset)
	}
	offset = 8
	if noffset = NeedleOffset(offset); noffset != uint32(offset/int64(PaddingSize)) {
		t.Errorf("noffset: %d not match", noffset)
		t.FailNow()
	}else{
		t.Log(noffset)
	}
}

func TestBlockOffset(t *testing.T) {
	var (
		offset  int64
		noffset uint32
	)
	noffset = 1
	if offset = BlockOffset(noffset); offset != int64(noffset*PaddingSize) {
		t.Errorf("offset: %d not match", offset)
		t.FailNow()
	}else{
		t.Log(offset)
	}
	noffset = 2
	if offset = BlockOffset(noffset); offset != int64(noffset*PaddingSize) {
		t.Errorf("offset: %d not match", offset)
		t.FailNow()
	}else{
		t.Log(offset)
	}
	noffset = 4
	if offset = BlockOffset(noffset); offset != int64(noffset*PaddingSize) {
		t.Errorf("offset: %d not match", offset)
		t.FailNow()
	}else{
		t.Log(offset)
	}
}

func TestTestNeedleSize(t *testing.T) {
	n:=Size(4)
	if n != 40 {
		t.FailNow()
	}
}

func compareNeedle(t *testing.T, n *Needle, key int64, cookie int32, data []byte, flag byte, checksum uint32) {
	if n.Key != key || n.Cookie != cookie || !bytes.Equal(n.Data, data) || n.Flag != flag || n.Checksum != checksum {
		t.Errorf("not match: %s, %d, %d, %d", n, key, cookie, checksum)
		t.FailNow()
	}
}
