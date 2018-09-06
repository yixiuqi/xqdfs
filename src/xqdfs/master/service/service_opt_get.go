package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/master/strategy"
	"xqdfs/utils/plugin"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdOptGet,ServiceOptGet)
}

/**
 * @api {post} /opt/get [Opt]图片下载
 * @apiDescription [Opt]图片下载
 * @apiGroup Master
 * @apiVersion 1.0.0
 * @apiParam {string} url 图片url
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (成功返回参数) {string} img 图片
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
func ServiceOptGet(m map[string]interface{}) interface{}{
	var strategyServer *strategy.AllocStrategyServer
	if s:=plugin.PluginGetObject(plugin.PluginStrategyServer);s==nil {
		log.Errorf("%s no support",plugin.PluginStrategyServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		strategyServer=s.(*strategy.AllocStrategyServer)
	}

	var url string
	value,ok:=m["url"]
	if ok {
		url=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"url missing")
	}

	img,err:=strategyServer.Read(url)
	if err!=nil {
		log.Error(err)
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
		}
	}else{
		json:=gabs.New()
		json.Set(img,"img")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
