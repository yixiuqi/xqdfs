package main

import (
	"flag"
	"xqdfs/master/conf"
	"xqdfs/master/server"
	"xqdfs/utils/log"
)

const(
	Ver	 = "1.0.0"
)

func main() {
	var (
		c	*conf.Config
		httpServer *server.HttpServer
		err	error
		configFile string
	)

	flag.StringVar(&configFile, "c", "./store.toml", " set master config file path")
	flag.Parse()
	log.Infof("xqdfs master[%s] start", Ver)

	if c, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFile, err)
		return
	}
	log.Info("config")
	log.Info(c)

	if httpServer, err = server.NewHttpServer(c); err != nil {
		log.Errorf("http server init error[%v]",err)
		return
	}

	log.Info("wait signal...")
	StartSignal(httpServer)
}
