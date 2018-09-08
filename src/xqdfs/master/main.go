package main

import (
	"flag"
	"fmt"
	"time"

	"xqdfs/proxy"
	"xqdfs/channel"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/master/conf"
	"xqdfs/utils/plugin"
	"xqdfs/master/strategy"
	_"xqdfs/master/service"
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
		server *channel.Server
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

	if server, err = channel.NewServer(config.Server); err != nil {
		log.Errorf("create server error[%v]",err)
		return
	}
	log.Info("system start")
	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(server,configureServer,discoveryServer,strategyServer,proxyStorage)
}

func logo(){
	time.Sleep(time.Millisecond*500)
	fmt.Println("")
	fmt.Println(" #     #    ###    #####    #######   #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("  #   #    #   #   #    #   #        #     # ");time.Sleep(time.Millisecond*100)
	fmt.Println("   # #    #     #  #     #  #        #       ");time.Sleep(time.Millisecond*100)
	fmt.Println("    #     #     #  #     #  #####     #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("   # #    #     #  #     #  #              # ");time.Sleep(time.Millisecond*100)
	fmt.Println("  #   #    #   #   #    #   #        #     # ");time.Sleep(time.Millisecond*100)
	fmt.Println(" #     #    ###    #####    #         #####  ");time.Sleep(time.Millisecond*100)
	fmt.Println("              ###                            ");time.Sleep(time.Millisecond*100)
}