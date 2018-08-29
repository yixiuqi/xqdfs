package helper

import (
	"net/http"
	"io/ioutil"
)

func URLDownload(url string) ([]byte,error){
	res, err := http.Get(url)
	if err != nil {
		return nil,err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil,err
	}

	return body,nil
}


