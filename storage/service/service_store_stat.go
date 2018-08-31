package service

import (
	"encoding/json"
	"bytes"

	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
	"xqdfs/errors"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.HttpStoreStat,ServiceStoreStat)
}

/**
 * @api {post} /store/stat [Store]查询状态
 * @apiDescription [Store]查询状态
 * @apiGroup Storage
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
func ServiceStoreStat(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	jsonStat:=gabs.New()
	jsonStat.Set(storage.FreeId,"freeId")

	jsonStat.Array("volumes")
	for _,v:=range storage.Volumes{
		b,err:=json.Marshal(v)
		if err==nil{
			dec := json.NewDecoder(bytes.NewBuffer(b))
			dec.UseNumber()
			item,err:=gabs.ParseJSONDecoder(dec)
			item.Set(v.ImageCount(),"imageCount")
			if err==nil {
				jsonStat.ArrayAppend(item.Data(), "volumes")
			}
		}else{
			log.Error(err)
		}
	}

	jsonStat.Array("freeVolumes")
	for _,v:=range storage.FreeVolumes{
		b,err:=json.Marshal(v)
		if err==nil{
			dec := json.NewDecoder(bytes.NewBuffer(b))
			dec.UseNumber()
			item,err:=gabs.ParseJSONDecoder(dec)
			if err==nil{
				jsonStat.ArrayAppend(item.Data(),"freeVolumes")
			}
		}else{
			log.Error(err)
		}
	}

	stat,err:=json.Marshal(storage.Stats)
	if err==nil{
		dec := json.NewDecoder(bytes.NewBuffer(stat))
		dec.UseNumber()
		item,err:=gabs.ParseJSONDecoder(dec)
		if err==nil{
			jsonStat.Set(item.Data(),"stats")
		}
	}else{
		log.Error(err)
	}

	result:=helper.ResultBuildWithBody(constant.Success,jsonStat)
	return result
}
