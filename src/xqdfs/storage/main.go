package main

import (
	"fmt"
	"time"
	"flag"

	"xqdfs/errors"
	"xqdfs/channel"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/storage/conf"
	"xqdfs/storage/store"
	"xqdfs/configure/defines"
	"xqdfs/storage/replication"
	_"xqdfs/storage/service"
)

const(
	Version	 = "1.0.0"
)

func main() {
	var (
		configFilePath string
		config	*conf.Config
		s	*store.Store
		server *channel.Server
		configureServer *configure.ConfigureServer
		replicationServer *replication.ReplicationServer
		err	error
	)

	flag.StringVar(&configFilePath, "-c", "./store.toml", " set store config file path")
	flag.Parse()
	log.Infof("xqdfs store version[%s] start", Version)

	if config, err = conf.NewConfig(configFilePath); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFilePath, err)
		return
	}
	plugin.PluginAddObject(plugin.PluginLocalConfig,config)

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
		plugin.PluginAddObject(plugin.PluginConfigure,configureServer)
	}

	if s, err = store.NewStore(config); err != nil {
		log.Errorf("store init error[%v]",err)
		return
	}else{
		err=s.Init()
		if err==errors.ErrVolumeExist{
			log.Info("store already init")
		}
		plugin.PluginAddObject(plugin.PlugineStorage,s)
	}

	if replicationServer, err = replication.NewReplicationServer(config,s,configureServer); err != nil {
		log.Errorf("create sync server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginReplicationServer,replicationServer)
	}

	if server, err = channel.NewServer(config.Server); err != nil {
		log.Errorf("create server error[%v]",err)
		return
	}
	log.Info("system start")
	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(configureServer,s,replicationServer,server)
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
