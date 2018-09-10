package main

import (
	"os/signal"
	"syscall"
	"os"

	"xqdfs/channel"
	"xqdfs/utils/log"
	"xqdfs/configure"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
)

func StartSignal(configureServer *configure.ConfigureServer,store *store.Store,replicationServer *replication.ReplicationServer,server *channel.Server) {
	var (
		c chan os.Signal
		s os.Signal
	)
	c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGINT, syscall.SIGSTOP)
	for {
		s = <-c
		log.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			server.Stop()
			replicationServer.Stop()
			store.Close()
			configureServer.Stop()
			log.Info("system stop")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
