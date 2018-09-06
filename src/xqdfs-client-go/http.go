package xqdfs_client_go

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"net"
	"time"
	"errors"
)

var(
	httpClient *http.Client
)

func init() {
	httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        16,
			MaxIdleConnsPerHost: 16,
			IdleConnTimeout:	 time.Duration(30)* time.Second,
		},
		Timeout: 30 * time.Second,
	}
	return client
}

func HttpPost(host string,sendData []byte)([]byte,error) {
	req, err := http.NewRequest("POST", host, bytes.NewBuffer(sendData))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
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
