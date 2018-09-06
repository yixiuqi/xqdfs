package helper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"errors"
)

func HttpPost(client *http.Client,host string,sendData []byte)([]byte,error) {
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(sendData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	if resp.Status=="200 OK"{
		return data,nil
	}else{
		return nil,errors.New(resp.Status)
	}
}
