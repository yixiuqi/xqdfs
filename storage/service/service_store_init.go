package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"
)

/**
 * @api {post} /store/init [Store]初始化
 * @apiDescription [Store]初始化
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/store/init
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
func ServiceStoreInit(context *Context,m map[string]interface{}) interface{}{
	err:=context.Store.Init()
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
			context.ReplicationServer.Replication(p)
		}else{
			log.Debug("receive replication request")
		}
		return helper.ResultBuild(constant.Success)
	}
}
