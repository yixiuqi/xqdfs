package binlog

import (
	"testing"
)

func TestBinlogWriter_Write(t *testing.T) {
	b,err:=NewBinlogWriter("1.log")
	if err!=nil{
		t.Error(err)
	}
	defer b.Close()

	b.Write(BinlogUpload,1,6592784829277647938,123456789)
}
