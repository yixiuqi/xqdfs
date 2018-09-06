package binlog

import (
	"testing"
	"fmt"
)

func TestNewBinlogConfigure(t *testing.T) {
	conf,err:=NewBinlogConfigure("./",1,"192.168.10.25:10086")
	if err!=nil{
		t.Error(err)
	}
	fmt.Println(conf)

	conf.StorageId=123
	conf.BinlogLine=1110
	conf.Update()
	fmt.Println(conf)

	conf.Close()
}
