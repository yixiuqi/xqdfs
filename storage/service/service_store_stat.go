package service

import (
	"xqdfs/utils/helper"
	"github.com/Jeffail/gabs"
	"encoding/json"
	"xqdfs/utils/log"
)

/**
 * @api {post} /store/stat 查询状态
 * @apiDescription 查询状态
 * @apiGroup Info
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/store/stat
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
func ServiceStoreStat(context *Context,m map[string]interface{}) interface{}{
	jsonStat:=gabs.New()
	jsonStat.Set(context.Store.FreeId,"freeId")

	jsonStat.Array("volumes")
	for _,v:=range context.Store.Volumes{
		b,err:=json.Marshal(v)
		if err==nil{
			item,err:=gabs.ParseJSON(b)
			if err==nil {
				jsonStat.ArrayAppend(item.Data(), "volumes")
			}
		}else{
			log.Error(err)
		}
	}

	jsonStat.Array("freeVolumes")
	for _,v:=range context.Store.FreeVolumes{
		b,err:=json.Marshal(v)
		if err==nil{
			item,err:=gabs.ParseJSON(b)
			if err==nil{
				jsonStat.ArrayAppend(item.Data(),"freeVolumes")
			}
		}else{
			log.Error(err)
		}
	}

	stat,err:=json.Marshal(context.Store.Stats)
	if err==nil{
		item,err:=gabs.ParseJSON(stat)
		if err==nil{
			jsonStat.Set(item.Data(),"stat")
		}
	}else{
		log.Error(err)
	}

	result:=helper.ResultBuildWithBody(Success,jsonStat)
	return result
}
