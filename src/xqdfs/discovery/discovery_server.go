package discovery

import (
	"xqdfs/configure"
	"xqdfs/discovery/http"
	"xqdfs/discovery/defines"
)

type DiscoveryServer struct {
	discovery defines.Discovery
}

func NewDiscoveryServer(conf *configure.ConfigureServer) (*DiscoveryServer,error){
	d,err:=http.NewDiscoveryHttp(conf)
	if err!=nil {
		return nil,err
	}

	s:=&DiscoveryServer{
		discovery:d,
	}

	return s,nil
}

func (this *DiscoveryServer) Groups() []*defines.Group {
	return this.discovery.Groups()
}

func (this *DiscoveryServer) Storages() []*defines.Storage {
	return this.discovery.Storages()
}

func (this *DiscoveryServer) Stop() {
	this.discovery.Stop()
}
