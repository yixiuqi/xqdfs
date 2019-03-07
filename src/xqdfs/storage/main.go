package main

import (
	"fmt"
	"time"
	"flag"

	"xqdfs/errors"
	"xqdfs/channel"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
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
		err	error
		startTime int64
		endTime int64
		configFilePath string
		config	*conf.Config
		storageSvr	*store.Store
		channelSvr *channel.Server
		confSvr *configure.ConfigureServer
		replicationServer *replication.ReplicationServer
	)

	flag.StringVar(&configFilePath, "-c", "./xqdfs_storage.toml", "storage config file path")
	flag.Parse()
	log.Infof("xqdfs storage version[%s] start", Version)

	if config, err = conf.NewConfig(configFilePath); err != nil {
		log.Errorf("NewConfig[%s] error[%v]", configFilePath, err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginLocalConfig,config)
	}

	if confSvr,err = configure.NewConfigureServer(config.Configure.Param); err != nil {
		log.Errorf("create configure server error[%v]",err)
		return
	}else{
		curStorageId:=int32(config.Server.Id)
		c,_:=confSvr.StorageGet(curStorageId)
		if c==nil {
			storageDal:=&defines.StorageDal{
				Id:curStorageId,
				Addr:fmt.Sprintf("%s:%d",config.Server.Host,config.Server.Port),
				Desc:config.Server.Desc,
			}
			err=confSvr.StorageAdd(storageDal)
			if err!=nil{
				log.Errorf("configure error[%v]",err)
				return
			}
		}
		plugin.PluginAddObject(plugin.PluginConfigure,confSvr)
	}

	startTime=helper.CurrentTime()
	if storageSvr, err = store.NewStore(config); err != nil {
		log.Errorf("NewStore error[%v]",err)
		return
	}else{
		err=storageSvr.Init()
		if err==errors.ErrVolumeExist {
			log.Info("storage already init")
		}
		plugin.PluginAddObject(plugin.PlugineStorage,storageSvr)
		endTime=helper.CurrentTime()
		log.Info("NewStore elapse:",endTime-startTime)
	}

	if replicationServer, err = replication.NewReplicationServer(config,storageSvr,confSvr); err != nil {
		log.Errorf("create sync server error[%v]",err)
		return
	}else{
		plugin.PluginAddObject(plugin.PluginReplicationServer,replicationServer)
	}

	if channelSvr, err = channel.NewServer(config.Server); err != nil {
		log.Errorf("create server error[%v]",err)
		return
	}

	log.SetLevel(config.Log.Level)
	go logo()
	StartSignal(confSvr,storageSvr,replicationServer,channelSvr)
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
