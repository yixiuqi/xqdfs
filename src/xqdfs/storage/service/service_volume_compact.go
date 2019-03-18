package service

import (
	"sync"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/storage/replication/process"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeCompact,ServiceVolumeCompact)
	plugin.PluginAddService(constant.CmdVolumeCompactStatus,ServiceVolumeCompactStatus)
}

var(
	compactTaskCount int32
	lock sync.RWMutex
)
/**
 * @api {post} /volume/compact [Volume]卷压缩
 * @apiDescription [Volume]卷压缩
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {int32} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
type RequestVolumeCompact struct {
	Vid int32 			`json:"vid"`
	Replication bool 	`json:"replication"`
}
func ServiceVolumeCompact(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeCompact{}
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

	if len(storage.FreeVolumes) == 0 {
		log.Error(errors.ErrStoreNoFreeVolume.Error())
		return helper.ResultBuildWithExtInfo(errors.RetStoreNoFreeVolume,errors.ErrStoreNoFreeVolume.Error())
	}

	if v:= storage.Volumes[req.Vid]; v != nil {
		if v.Compact {
			return helper.ResultBuildWithExtInfo(errors.RetVolumeInCompact,errors.ErrVolumeInCompact.Error())
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}

	lock.Lock()
	if compactTaskCount>1{
		lock.Unlock()
		return helper.ResultBuildWithExtInfo(errors.RetVolumeTooManyCompact,errors.ErrVolumeTooManyCompact.Error())
	}else{
		compactTaskCount++
		lock.Unlock()
	}

	go func(replication bool,vid int32) {
		defer helper.HandleErr()

		if replication==true {
			p:=&process.ReplicationStorageVolumeCompact{
				Vid:vid,
			}
			replicationServer.Replication(p)
		}

		log.Infof(">>> start compact volume[%v]",vid)
		start:=helper.CurrentTime()
		err:=storage.CompactVolume(vid)
		end:=helper.CurrentTime()
		log.Infof("<<< finish volume[%v] compact,elapse[%d]",vid,end-start)

		if err!=nil{
			log.Infof("volume[%v] compact error",err)
		}else{
			log.Infof("volume[%v] compact ok",vid)
		}

		lock.Lock()
		compactTaskCount--
		lock.Unlock()
	}(req.Replication,req.Vid)
	return helper.ResultBuild(constant.Success)
}

type RequestVolumeCompactStatus struct {
	Vid int32 `json:"vid"`
}
func ServiceVolumeCompactStatus(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeCompactStatus{}
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

	v:= storage.Volumes[req.Vid]
	if v != nil {
		json:=gabs.New()
		json.Set(v.IsCompact(),"status")
		json.Set(compactTaskCount,"task_count")
		return helper.ResultBuildWithBody(constant.Success,json)
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}

	return helper.ResultBuild(constant.Success)
}

