package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/strategy"

	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdOptDelete,ServiceOptDelete)
}

/**
 * @api {post} /opt/delete [Opt]图片删除
 * @apiDescription [Opt]图片删除
 * @apiGroup Master
 * @apiVersion 1.0.0
 * @apiParam {string} url 图片url
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (成功返回参数) {string} img 图片
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
type RequestOptDelete struct {
	Url string 			`json:"url"`
}
func ServiceOptDelete(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestOptDelete{}
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

	err=strategyServer.Delete(req.Url)
	if err!=nil {
		log.Warn(err)
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
