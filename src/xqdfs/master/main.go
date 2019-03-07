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
	Version	 = "1.0.0"
)

func main() {
	var (
		err	error
		configFilePath string
		config	*conf.Config
		proxyStorage *proxy.ProxyStorage
		confSvr *configure.ConfigureServer
		discoverySvr *discovery.DiscoveryServer
		strategySvr *strategy.AllocStrategyServer
		channelSvr *channel.Server
	)

	flag.StringVar(&configFilePath, "-c", "./xqdfs_master.toml", " set master config file path")
	flag.Parse()
	log.Infof("xqdfs master version[%s] start", Version)

	if config, err = conf.NewConfig(configFilePath); err != nil {
		log.Errorf("NewConfig[%s] error[%v]", configFilePath, err)
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

	if confSvr,err = configure.NewConfigureServer(config.Configure.Param); err != nil {
		log.Errorf("create configure server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginConfigure,confSvr)
	}

	if discoverySvr,err = discovery.NewDiscoveryServer(confSvr); err != nil {
		log.Errorf("create discovery server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginDiscoveryServer,discoverySvr)
	}

	if strategySvr,err = strategy.NewAllocStrategyServer(); err != nil {
		log.Errorf("create alloc strategy error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginStrategyServer,strategySvr)
	}

	if channelSvr, err = channel.NewServer(config.Server); err != nil {
		log.Errorf("create server error[%v]",err)
		return
	}

	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(channelSvr,confSvr,discoverySvr,strategySvr,proxyStorage)
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