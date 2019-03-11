package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/strategy"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
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
type RequestOptGet struct {
	Url string 			`json:"url"`
}
func ServiceOptGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestOptGet{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var strategyServer *strategy.AllocStrategyServer
	if s:=plugin.PluginGetObject(plugin.PluginStrategyServer);s==nil {
		log.Errorf("%s no support",plugin.PluginStrategyServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		strategyServer=s.(*strategy.AllocStrategyServer)
	}

	img,err:=strategyServer.Read(req.Url)
	if err!=nil {
		log.Warn(err)
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
