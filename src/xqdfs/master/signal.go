package main

import (
	"os/signal"
	"syscall"
	"os"

	"xqdfs/proxy"
	"xqdfs/channel"
	"xqdfs/utils/log"
	"xqdfs/master/configure"
	"xqdfs/discovery"
	"xqdfs/master/strategy"
)

func StartSignal(server *channel.Server,
								configureServer *configure.ConfigureServer,
								discoveryServer *discovery.DiscoveryServer,
								strategyServer *strategy.AllocStrategyServer,
								proxyStorage *proxy.ProxyStorage) {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGINT, syscall.SIGSTOP)
	// Block until a signal is received.
	for {
		s = <-c
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			strategyServer.Stop()
			proxyStorage.Stop()
			discoveryServer.Stop()
			configureServer.Stop()
			server.Stop()
			log.Info("system stop")
			return
		case syscall.SIGHUP:
			// TODO reload
			//return
		default:
			return
		}
	}
}
