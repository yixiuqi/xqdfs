package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/storage/replication/process"

	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeDelete,ServiceVolumeDelete)
}

type RequestVolumeDelete struct {
	Vid int32 			`json:"vid"`
	Key int64 			`json:"key"`
	Replication bool 	`json:"replication"`
}
func ServiceVolumeDelete(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeDelete{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Warnf("%s no support",plugin.PlugineStorage)
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

	v:= storage.Volumes[req.Vid]
	if v != nil {
		err:= v.Delete(req.Key)
		if req.Replication==true {
			p:=&process.ReplicationDelete{
				Vid:req.Vid,
				Key:req.Key,
			}
			replicationServer.Replication(p)
		}
		if err!=nil{
			//log.Error(err)
			e,ok:=err.(errors.Error)
			if ok {
				return helper.ResultBuildWithExtInfo(int32(e),err.Error())
			}else{
				return helper.ResultBuildWithExtInfo(errors.RetOptDelete,err.Error())
			}
		}else{
			return helper.ResultBuild(constant.Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
