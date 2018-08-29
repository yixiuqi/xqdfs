package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"
)

/**
 * @api {post} /volume/clear [Volume]块清空
 * @apiDescription [Volume]块清空
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/clear
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
func ServiceVolumeClear(context *Context,m map[string]interface{}) interface{}{
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
		err:= v.Clear()
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetVolumeClear,err.Error())
		}else{
			replication,ok:=m["replication"]
			if ok && replication==true {
				log.Debug("ReplicationStorageVolumeClear")
				p:=&process.ReplicationStorageVolumeClear{
					Vid:vid,
				}
				context.ReplicationServer.Replication(p)
			}else{
				log.Debug("receive replication request")
			}
			return helper.ResultBuild(constant.Success)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
