package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"

	"github.com/Jeffail/gabs"
)

/**
 * @api {post} /volume/addfree [Volume]块创建
 * @apiDescription [Volume]块创建
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int} count 创建数目
 * @apiParam {string} bdir block文件目录
 * @apiParam {string} idir block.idx文件目录
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/addfree
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
func ServiceVolumeAddFree(context *Context,m map[string]interface{}) interface{}{
	var count int = 1
	var bdir string
	var idir string

	value,ok:=m["count"]
	if ok {
		tmp,err:=helper.GetInt(value)
		if err==nil{
			count=int(tmp)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"count missing")
	}

	value,ok=m["bdir"]
	if ok {
		bdir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"bdir missing")
	}

	value,ok=m["idir"]
	if ok {
		idir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"idir missing")
	}

	sn,err:=context.Store.AddFreeVolume(count,bdir,idir)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetVolumeAddFree,err.Error())
	}else{
		json:=gabs.New()
		json.Set(sn,"count")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
