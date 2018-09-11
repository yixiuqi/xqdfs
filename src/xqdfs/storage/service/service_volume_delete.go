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
	plugin.PluginAddService(constant.CmdVolumeDelete,ServiceVolumeDelete)
}

func ServiceVolumeDelete(m map[string]interface{}) interface{}{
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

	var vid int32
	var key int64
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"vid missing")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"key missing")
	}

	v:= storage.Volumes[vid]
	if v != nil {
		err:= v.Delete(key)
		replication,ok:=m["replication"]
		if ok && replication==true {
			p:=&process.ReplicationDelete{
				Vid:vid,
				Key:key,
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
