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
type RequestVolumeClear struct {
	Vid int32 			`json:"vid"`
	Replication bool 	`json:"replication"`
}
func ServiceVolumeClear(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeClear{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

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

	v:= storage.Volumes[req.Vid]
	if v != nil {
		err:= v.Clear()
		if err!=nil{
			log.Warn(err)
			return helper.ResultBuildWithExtInfo(errors.RetVolumeClear,err.Error())
		}else{
			if req.Replication==true {
				p:=&process.ReplicationStorageVolumeClear{
					Vid:req.Vid,
				}
				replicationServer.Replication(p)
			}
			return helper.ResultBuild(constant.Success)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
