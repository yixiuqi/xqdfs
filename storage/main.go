package main

import (
	"flag"
	"fmt"
	"time"

	"xqdfs/utils/log"
	"xqdfs/storage/conf"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/storage/service"
	"xqdfs/errors"
	"xqdfs/configure"
	"xqdfs/configure/defines"
)

const(
	Ver	 = "1.0.0"
)

func main() {
	var (
		configFile string
		config	*conf.Config
		s	*store.Store
		httpServer *service.HttpServer
		configureServer *configure.ConfigureServer
		replicationServer *replication.ReplicationServer
		err	error
	)

	flag.StringVar(&configFile, "c", "./store.toml", " set store config file path")
	flag.Parse()
	log.Infof("xqdfs store[%s] start", Ver)

	if config, err = conf.NewConfig(configFile); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFile, err)
		return
	}

	if configureServer,err = configure.NewConfigureServer(config.Configure.Param); err != nil {
		log.Errorf("create configure server error[%v]",err)
		return
	}else{
		id:=int32(config.Server.Id)
		c,_:=configureServer.StorageGet(id)
		if c==nil {
			storageDal:=&defines.StorageDal{
				Id:id,
				Addr:fmt.Sprintf("%s:%d",config.Server.Host,config.Server.Port),
				Desc:config.Server.Desc,
			}
			err=configureServer.StorageAdd(storageDal)
			if err!=nil{
				log.Errorf("configure error[%v]",err)
				return
			}
		}
	}

	if s, err = store.NewStore(config); err != nil {
		log.Errorf("store init error[%v]",err)
		return
	}else{
		err=s.Init()
		if err==errors.ErrVolumeExist{
			log.Info("store already init")
		}
	}

	if replicationServer, err = replication.NewReplicationServer(config,s,configureServer); err != nil {
		log.Errorf("sync server init error[%v]",err)
		return
	}

	if httpServer, err = service.NewHttpServer(config,s,replicationServer); err != nil {
		log.Errorf("http server init error[%v]",err)
		return
	}

	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(configureServer,s,replicationServer,httpServer)
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
