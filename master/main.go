package main

import (
	"flag"
	"fmt"
	"time"

	_"xqdfs/master/service"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/proxy"
	"xqdfs/configure"
	"xqdfs/master/conf"
	"xqdfs/master/strategy"
	"xqdfs/utils/plugin"
	"xqdfs/channel"
)

const(
	Ver	 = "1.0.0"
)

func main() {
	var (
		err	error
		configFile string
		config	*conf.Config
		proxyStorage *proxy.ProxyStorage
		configureServer *configure.ConfigureServer
		discoveryServer *discovery.DiscoveryServer
		strategyServer *strategy.AllocStrategyServer
		httpServer *channel.HttpServer
	)

	flag.StringVar(&configFile, "c", "./store.toml", " set master config file path")
	flag.Parse()
	log.Infof("xqdfs master[%s] start", Ver)

	if config, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("create config(\"%s\") error(%v)", configFile, err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginLocalConfig,config)
	}

	if proxyStorage,err = proxy.NewProxyStorage(); err != nil {
		log.Errorf("create proxy error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginProxyStorage,proxyStorage)
	}

	if configureServer,err = configure.NewConfigureServer(config.Configure.Param); err != nil {
		log.Errorf("create configure server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginConfigure,configureServer)
	}

	if discoveryServer,err = discovery.NewDiscoveryServer(configureServer); err != nil {
		log.Errorf("create discovery server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginDiscoveryServer,discoveryServer)
	}

	if strategyServer,err = strategy.NewAllocStrategyServer(); err != nil {
		log.Errorf("create alloc strategy error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginStrategyServer,strategyServer)
	}

	if httpServer, err = channel.NewHttpServer(config.Server.Port); err != nil {
		log.Errorf("create http server error[%v]",err)
		return
	}

	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(httpServer,configureServer,discoveryServer,strategyServer,proxyStorage)
}

func logo(){
	time.Sleep(time.Millisecond*500)
	fmt.Println(" #     #    ###    #####    #######   #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("  #   #    #   #   #    #   #        #     # ");time.Sleep(time.Millisecond*100)
	fmt.Println("   # #    #     #  #     #  #        #       ");time.Sleep(time.Millisecond*100)
	fmt.Println("    #     #     #  #     #  #####     #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("   # #    #     #  #     #  #              # ");time.Sleep(time.Millisecond*100)
	fmt.Println("  #   #    #   #   #    #   #        #     # ");time.Sleep(time.Millisecond*100)
	fmt.Println(" #     #    ###    #####    #         #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("              ###                            ");time.Sleep(time.Millisecond*100)
}