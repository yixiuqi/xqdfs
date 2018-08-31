package strategy

import (
	"xqdfs/proxy"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/master/conf"
	"xqdfs/master/strategy/alloc"
	"xqdfs/master/strategy/clear"
	"xqdfs/master/strategy/defines"
	"xqdfs/utils/log"
)

type AllocStrategyServer struct {
	allocStrategy defines.AllocStrategy
	clearStrategy defines.ClearStrategy
}

func NewAllocStrategyServer(conf *conf.Config,configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) (*AllocStrategyServer,error){
	clear,err:=clear.NewClearTimeOld(configureServer,discoveryServer,proxyStorage)
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	alloc,err:=order.NewAllocOrder(configureServer,discoveryServer,proxyStorage)
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	s:=&AllocStrategyServer{
		allocStrategy:alloc,
		clearStrategy:clear,
	}

	return s,nil
}

func (this *AllocStrategyServer) Write(key int64,cookie int32,img []byte) (string,error) {
	return this.allocStrategy.Write(key,cookie,img)
}

func (this *AllocStrategyServer) Read(url string) ([]byte,error) {
	return this.allocStrategy.Read(url)
}

func (this *AllocStrategyServer) Delete(url string) error {
	return this.allocStrategy.Delete(url)
}

func (this *AllocStrategyServer) Stop() {
	this.allocStrategy.Stop()
	this.clearStrategy.Stop()
}
