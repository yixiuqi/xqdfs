package xqdfs_client_go

import (
	"testing"
	"fmt"
)

func TestClientUpload(t *testing.T) {
	img:=make([]byte,1024)
	url,err:=ClientUpload("http://192.168.10.25:10087/opt/upload",img)
	if err!=nil {
		t.Error(err)
	}else{
		fmt.Println(url)
	}
}
