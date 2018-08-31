package service

import (
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/utils/plugin"
	"xqdfs/storage/conf"
	"xqdfs/utils/helper"
)

func init() {
	plugin.PluginAddService(constant.HttpStoreConf,ServiceStoreConf)
}

/**
 * @api {post} /store/conf [Store]查询配置
 * @apiDescription [Store]查询配置
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/store/conf
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
func ServiceStoreConf(m map[string]interface{}) interface{}{
	var config *conf.Config
	if c:=plugin.PluginGetObject(plugin.PluginLocalConfig);c==nil {
		log.Errorf("%s no support",plugin.PluginLocalConfig)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		config=c.(*conf.Config)
	}

	json,err:=config.Json()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStoreConfigure,err.Error())
	}else{
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
