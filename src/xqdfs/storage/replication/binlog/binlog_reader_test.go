package binlog

import (
	"testing"
	"fmt"
)

func TestBinlogReader_Read(t *testing.T) {
	b,err:=NewBinlogReader("1.log")
	if err!=nil{
		t.Error(err)
	}
	defer b.Close()

	binlog,err:=b.Read(5)
	if err!=nil {
		t.Error(err)
	}
	fmt.Print(binlog)
}
