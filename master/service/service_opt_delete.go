package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
)

/**
 * @api {post} /opt/delete [Opt]图片删除
 * @apiDescription [Opt]图片删除
 * @apiGroup Master
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {string} url 图片url
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (成功返回参数) {string} img 图片
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/opt/delete
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
func ServiceOptDelete(context *Context,m map[string]interface{}) interface{}{
	var url string

	value,ok:=m["url"]
	if ok {
		url=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"url missing")
	}

	err:=context.StrategyServer.Delete(url)
	if err!=nil {
		log.Error(err)
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetOptDelete,err.Error())
		}
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
