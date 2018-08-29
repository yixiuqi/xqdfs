package ssdb

import (
	"testing"
	"fmt"

	"xqdfs/utils/log"
	"xqdfs/utils/helper"
)

func TestHSet(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	hset:=NewSSDBHash()
	hset.HSet("h1","a1","1111")
}

func TestHGet(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	hget:=NewSSDBHash()
	v,err:=hget.HGet("h1","a1")
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(v)
	}
}

func TestHGetAll(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	hget:=NewSSDBHash()
	v,err:=hget.HGetAll("h1")
	if err!=nil{
		fmt.Println(err)
	}else{
		fmt.Println(v)
	}
}

func TestHash(t *testing.T) {
	log.SetLevel("error")
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	max:=1000000
	for i:=0;i<max;i++{
		k:=helper.IntToString(i)
		hget:=NewSSDBHash()
		err:=hget.HSet("test1",k,helper.UUIDBuild())
		if err!=nil{
			fmt.Println(err)
		}

		if (i%1000)==0{
			fmt.Println("HSet:",i)
		}
	}

	for i:=0;i<max;i++{
		k:=helper.IntToString(i)
		hget:=NewSSDBHash()
		err:=hget.HDel("test1",k)
		if err!=nil{
			fmt.Println(err)
		}

		if (i%1000)==0{
			fmt.Println("HDel:",i)
		}
	}
}
