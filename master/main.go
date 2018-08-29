package main

import (
	"flag"
	"fmt"
	"time"

	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/proxy"
	"xqdfs/configure"
	"xqdfs/master/conf"
	"xqdfs/master/strategy"
	"xqdfs/master/service"
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
		httpServer *service.HttpServer
	)

	flag.StringVar(&configFile, "c", "./store.toml", " set master config file path")
	flag.Parse()
	log.Infof("xqdfs master[%s] start", Ver)

	if config, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("create config(\"%s\") error(%v)", configFile, err)
		return
	}

	if proxyStorage,err = proxy.NewProxyStorage(); err != nil {
		log.Errorf("create proxy error[%v]",err)
		return
	}

	if configureServer,err = configure.NewConfigureServer(config.Configure.Param); err != nil {
		log.Errorf("create configure server error[%v]",err)
		return
	}

	if discoveryServer,err = discovery.NewDiscoveryServer(configureServer); err != nil {
		log.Errorf("create discovery server error[%v]",err)
		return
	}

	if strategyServer,err = strategy.NewAllocStrategyServer(config,discoveryServer,proxyStorage); err != nil {
		log.Errorf("create alloc strategy error[%v]",err)
		return
	}

	if httpServer, err = service.NewHttpServer(config,configureServer,discoveryServer,strategyServer,proxyStorage); err != nil {
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