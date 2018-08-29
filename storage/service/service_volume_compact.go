package service

import (
	"sync"
	
	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"

	"github.com/Jeffail/gabs"
)

var(
	compactTaskCount int32
	lock sync.RWMutex
)
/**
 * @api {post} /volume/compact [Volume]块回收
 * @apiDescription [Volume]块回收
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/compact
{
	"seq": "f428b43e-5bc1-4bea-a93a-ed177ca4552e"
}
 * @apiErrorExample 失败返回示例
{
    "info": "查询错误",
    "result": 1
}

 * @apiSuccessExample 成功返回示例
{
    "info": "success",
    "result": 0
}
* */
func ServiceVolumeCompact(context *Context,m map[string]interface{}) interface{}{
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

	if len(context.Store.FreeVolumes) == 0 {
		return helper.ResultBuildWithExtInfo(errors.RetStoreNoFreeVolume,errors.ErrStoreNoFreeVolume.Error())
	}

	if v:= context.Store.Volumes[vid]; v != nil {
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
			context.ReplicationServer.Replication(p)
		}else{
			log.Debug("receive replication request")
		}

		defer helper.HandleErr()
		start:=helper.CurrentTime()
		err:=context.Store.CompactVolume(vid)
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
 * @api {post} /volume/compact/status [Volume]块回收状态查询
 * @apiDescription [Volume]块回收状态查询
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/compact/status
{
	"seq": "f428b43e-5bc1-4bea-a93a-ed177ca4552e"
}
 * @apiErrorExample 失败返回示例
{
    "info": "查询错误",
    "result": 1
}

 * @apiSuccessExample 成功返回示例
{
    "info": "success",
    "result": 0
}
* */
func ServiceVolumeCompactStatus(context *Context,m map[string]interface{}) interface{}{
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

	v:= context.Store.Volumes[vid]
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

