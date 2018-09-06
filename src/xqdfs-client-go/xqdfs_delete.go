package xqdfs_client_go

import (
	"encoding/json"
	"errors"
)

type DeleteCmd struct {
	Url string 		`json:"url"`
}

type DeleteResult struct {
	Result int32 	`json:"result"`
	Info string 	`json:"info"`
}

func ClientDelete(host string,url string) (err error) {
	d:=&GetCmd{
		Url:url,
	}

	var data []byte
	data,err=json.Marshal(d)
	if err !=nil {
		return
	}

	var ret []byte
	ret,err=HttpPost(host,data)
	if err !=nil {
		return
	}

	r:=&DeleteResult{}
	err=json.Unmarshal(ret,r)
	if err !=nil {
		return
	}

	if r.Result!=0 {
		err=errors.New(r.Info)
	}
	return
}