package ssdb

import (
	"testing"
	"fmt"
)

func TestQSize(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	list:=NewSSDBList()
	size,err:=list.QSize("q1")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("size:",size)
	}
}

func TestQFront(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	list:=NewSSDBList()
	item,err:=list.QFront("q1")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("QFront:",item)
	}
}

func TestQPop(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	list:=NewSSDBList()
	item,err:=list.QPop("q1")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("QPop:",item)
	}
}

func TestQPush(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	list:=NewSSDBList()
	err:=list.QPush("q1","2")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("QPush ok")
	}
}

func TestQClear(t *testing.T) {
	SSDBConnectMgrInstance().Init("192.168.10.25:18888,192.168.10.31:18888")

	list:=NewSSDBList()
	err:=list.QClear("q1")
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("QClear ok")
	}
}
