package http

import (
	"net"
	"time"
	"net/http"
	"encoding/json"

	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/discovery/defines"
	
	"github.com/Jeffail/gabs"
)

var httpClient *http.Client
func init() {
	httpClient = createHTTPClient()
}

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 3 * time.Second,
			}).DialContext,
			MaxIdleConns:        2,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:	 time.Duration(10)* time.Second,
		},

		Timeout: 3 * time.Second,
	}
	return client
}

func HeartBeat(addr string) (storage *defines.Storage,err error) {
	url:="http://" + addr + "/store/stat"
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	var data []byte
	data,err=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warn(err)
		return
	}

	storage=&defines.Storage{}
	err=json.Unmarshal(data,storage)
	if err!=nil {
		log.Warn(err)
	}
	return
}