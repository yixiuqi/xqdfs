package http

import (
	"net"
	"time"
	"net/http"
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
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			MaxIdleConns:        16,
			MaxIdleConnsPerHost: 16,
			IdleConnTimeout:	 time.Duration(10)* time.Second,
		},

		Timeout: 10 * time.Second,
	}
	return client
}
