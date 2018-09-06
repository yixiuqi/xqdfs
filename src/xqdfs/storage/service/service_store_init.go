package service

import (
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/storage/replication/process"
)

func init() {
	plugin.PluginAddService(constant.CmdStoreInit,ServiceStoreInit)
}

func ServiceStoreInit(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var replicationServer *replication.ReplicationServer
	if r:=plugin.PluginGetObject(plugin.PluginReplicationServer);r==nil {
		log.Errorf("%s no support",plugin.PluginReplicationServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		replicationServer=r.(*replication.ReplicationServer)
	}

	err:=storage.Init()
	if err!=nil{
		log.Error(err)
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetStoreInitFailed,err.Error())
		}
	}else{
		replication,ok:=m["replication"]
		if ok && replication==true {
			p:=&process.ReplicationStorageInit{}
			replicationServer.Replication(p)
		}else{
			log.Debug("receive replication request")
		}
		return helper.ResultBuild(constant.Success)
	}
}
