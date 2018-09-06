package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
	"xqdfs/storage/replication"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeClear,ServiceVolumeClear)
}

/**
 * @api {post} /volume/clear [Volume]卷回收
 * @apiDescription [Volume]卷回收
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {int} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
func ServiceVolumeClear(m map[string]interface{}) interface{}{
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
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"vid missing")
	}

	v:= storage.Volumes[vid]
	if v != nil {
		err:= v.Clear()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetVolumeClear,err.Error())
		}else{
			replication,ok:=m["replication"]
			if ok && replication==true {
				p:=&process.ReplicationStorageVolumeClear{
					Vid:vid,
				}
				replicationServer.Replication(p)
			}else{
				log.Debug("receive replication request")
			}
			return helper.ResultBuild(constant.Success)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
