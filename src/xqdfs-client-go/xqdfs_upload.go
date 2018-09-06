package xqdfs_client_go

import (
	"encoding/json"
	"errors"
)

type UploadCmd struct {
	Img []byte		`json:"img"`
}

type UploadResult struct {
	Result int32 	`json:"result"`
	Info string 	`json:"info"`
	Url string 		`json:"url"`
}

func ClientUpload(host string,img []byte) (url string, err error) {
	u:=&UploadCmd{
		Img:img,
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

	r:=&UploadResult{}
	err=json.Unmarshal(ret,r)
	if err !=nil {
		return
	}

	if r.Result==0 {
		url=r.Url
	}else{
		err=errors.New(r.Info)
	}
	return
}
