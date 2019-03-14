package proxy

import (
	"xqdfs/utils/log"
	"xqdfs/proxy/http"
	"xqdfs/proxy/defines"

	"github.com/Jeffail/gabs"
)

type ProxyStorage struct {
	proxy defines.Proxy
}

func NewProxyStorage() (*ProxyStorage,error) {
	p:=&ProxyStorage{
		proxy:http.NewProxyHttp(),
	}
	return p,nil
}

func (this *ProxyStorage) Upload(host string,body *gabs.Container,vid int32,key int64,cookie int32,img []byte,replication bool) error {
	return this.proxy.Upload(host,body,vid,key,cookie,img,replication)
}

func (this *ProxyStorage) Get(host string,vid int32,key int64,cookie int32) ([]byte,error) {
	return this.proxy.Get(host,vid,key,cookie)
}

func (this *ProxyStorage) Delete(host string,vid int32,key int64,replication bool) error {
	return this.proxy.Delete(host,vid,key,replication)
}

func (this *ProxyStorage) StorageInit(host string,replication bool) error {
	return this.proxy.StorageInit(host,replication)
}

func (this *ProxyStorage) StorageVolumeCompact(host string,vid int32,replication bool) error {
	return this.proxy.StorageVolumeCompact(host,vid,replication)
}

func (this *ProxyStorage) StorageVolumeClear(host string,vid int32,replication bool) error {
	return this.proxy.StorageVolumeClear(host,vid,replication)
}

func (this *ProxyStorage) StorageGetConfigure(host string) (*gabs.Container,error) {
	return this.proxy.StorageGetConfigure(host)
}

func (this *ProxyStorage) Stop() {
	log.Info("ProxyStorage stop->")
	this.proxy.Stop()
	log.Info("ProxyStorage stop-<")
}