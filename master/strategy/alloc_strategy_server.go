package strategy

import (
	"xqdfs/master/conf"
	"xqdfs/master/strategy/order"
	strategydef "xqdfs/master/strategy/defines"
	"xqdfs/discovery"
	"xqdfs/proxy"
)

type AllocStrategyServer struct {
	allocStrategy strategydef.AllocStrategy
}

func NewAllocStrategyServer(conf *conf.Config,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) (*AllocStrategyServer,error){
	s:=&AllocStrategyServer{
		allocStrategy:order.NewAllocStrategyOrder(conf,discoveryServer,proxyStorage),
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
}
