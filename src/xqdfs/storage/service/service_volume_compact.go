package service

import (
	"sync"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/storage/replication/process"

	"github.com/Jeffail/gabs"
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
func ServiceVolumeCompact(m map[string]interface{}) interface{}{
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

	if len(storage.FreeVolumes) == 0 {
		log.Error(errors.ErrStoreNoFreeVolume.Error())
		return helper.ResultBuildWithExtInfo(errors.RetStoreNoFreeVolume,errors.ErrStoreNoFreeVolume.Error())
	}

	if v:= storage.Volumes[vid]; v != nil {
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

	replication,_:=m["replication"]
	go func(replication bool,vid int32) {
		if replication==true {
			p:=&process.ReplicationStorageVolumeCompact{
				Vid:vid,
			}
			replicationServer.Replication(p)
		}else{
			log.Debug("receive replication request")
		}

		defer helper.HandleErr()
		start:=helper.CurrentTime()
		err:=storage.CompactVolume(vid)
		end:=helper.CurrentTime()
		log.Debugf("volume[%v] compact time[%d]",vid,end-start)

		if err!=nil{
			log.Infof("volume[%v] compact error",err)
		}else{
			log.Infof("volume[%v] compact ok",vid)
		}

		lock.Lock()
		compactTaskCount--
		lock.Unlock()
	}(replication.(bool),vid)
	return helper.ResultBuild(constant.Success)
}

/**
 * @api {post} /volume/compact/status [Volume]卷压缩状态查询
 * @apiDescription [Volume]卷压缩状态查询
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {int32} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
func ServiceVolumeCompactStatus(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
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
		json:=gabs.New()
		json.Set(v.IsCompact(),"status")
		json.Set(compactTaskCount,"task_count")
		return helper.ResultBuildWithBody(constant.Success,json)
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}

	return helper.ResultBuild(constant.Success)
}
