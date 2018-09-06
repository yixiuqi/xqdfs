package ssdb

import (
	"testing"
	"fmt"
	"time"
)

func TestKVGet(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	for {
		kv:=NewSSDBKV()
		val,err:=kv.Get("a1")
		if err!=nil{
			fmt.Println(err)
		}else{
			fmt.Println("val:",len(val)," ",time.Now().String())
		}

		time.Sleep(time.Second)
	}
}

func TestKVSet(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	kv:=NewSSDBKV()
	v:=make([]byte,1024*1024)
	err:=kv.Set("a1",string(v))
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println("KVSet ok")
	}
}