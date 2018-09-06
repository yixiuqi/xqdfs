package xqdfs_client_go

import (
	"encoding/json"
	"errors"
)

type GetCmd struct {
	Url string 		`json:"url"`
}

type GetResult struct {
	Result int32 	`json:"result"`
	Info string 	`json:"info"`
	Img []byte		`json:"img"`
}

func ClientGet(host string,url string) (img []byte, err error) {
	u:=&GetCmd{
		Url:url,
	}

	var data []byte
	data,err=json.Marshal(u)
	if err !=nil {
		return
	}

	var ret []byte
	ret,err=HttpPost(host,data)
	if err !=nil {
		return
	}

	r:=&GetResult{}
	err=json.Unmarshal(ret,r)
	if err !=nil {
		return
	}

	if r.Result==0 {
		img=r.Img
	}else{
		err=errors.New(r.Info)
	}
	return
}
