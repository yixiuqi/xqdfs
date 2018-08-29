package main

import (
	"flag"
	"xqdfs/storage/conf"
	"xqdfs/utils/log"
	"xqdfs/storage/store"
	"xqdfs/utils/helper"
	"fmt"
	"xqdfs/storage/server"
)

const(
	Ver	 = "1.0.0"
)

func main() {
	var (
		c	*conf.Config
		s	*store.Store
		httpServer *server.HttpServer
		err	error
		configFile string
	)

	flag.StringVar(&configFile, "c", "./store.toml", " set store config file path")
	flag.Parse()
	log.Infof("xqdfs store[%s] start", Ver)

	if c, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFile, err)
		return
	}
	log.Info("config")
	log.Info(c)

	if s, err = store.NewStore(c); err != nil {
		log.Errorf("store init error[%v]",err)
		return
	}

	if httpServer, err = server.NewHttpServer(c,s); err != nil {
		log.Errorf("http server init error[%v]",err)
		return
	}

	log.Info("wait signal...")
	StartSignal(s,httpServer)
}

func test(){
	start:=helper.CurrentTime()
	items:=make(map[int64]string)
	for i:=int64(0);i<100000000;i++{
		id:=helper.KeyGenerate()
		items[id]=fmt.Sprintf("%v",i)
		if (i%100000)==0{
			fmt.Println(i)
		}

		if i==40000000{
			items[10086]="fdafdaf"
		}
	}
	end:=helper.CurrentTime()
	fmt.Println("all time:",end-start)

	for {
		start = helper.CurrentTime()
		for i := int64(0); i < 100000; i++ {
			id := helper.KeyGenerate()
			v, ok := items[id]
			if ok {
				v = v
			}
		}
		end = helper.CurrentTime()
		fmt.Println("time:", end-start)
	}
}
