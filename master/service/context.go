package service

import (
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/proxy"
	"xqdfs/master/conf"
	"xqdfs/master/strategy"
)

type Context struct {
	Conf *conf.Config
	ConfigureServer *configure.ConfigureServer
	DiscoveryServer *discovery.DiscoveryServer
	StrategyServer *strategy.AllocStrategyServer
	ProxyStorage *proxy.ProxyStorage
}

func NewContext(conf *conf.Config,configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,strategyServer *strategy.AllocStrategyServer,proxyStorage *proxy.ProxyStorage) *Context {
	context:=&Context{
		Conf:conf,
		ConfigureServer:configureServer,
		DiscoveryServer:discoveryServer,
		StrategyServer:strategyServer,
		ProxyStorage:proxyStorage,
	}
	return context
}
