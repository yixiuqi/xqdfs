package service

import (
	"context"
	"strings"
	"encoding/base64"

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
	plugin.PluginAddService(constant.CmdOptUpload,ServiceOptUpload)
}

/**
 * @api {post} /opt/upload [Opt]图片上传
 * @apiDescription [Opt]图片上传
 * @apiGroup Master
 * @apiVersion 1.0.0
 * @apiParam {string} img 图片数据
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
type RequestOptUpload struct {
	Img []byte `json:"img"`
}
func ServiceOptUpload(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestOptUpload{}
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

	key:=helper.KeyGenerate()
	url,err:=strategyServer.Write(key,constant.Cookie,req.Img)
	for err==errors.ErrNeedleExist {
		log.Error(err," try again")
		key=helper.KeyGenerate()
		url,err=strategyServer.Write(key,constant.Cookie,req.Img)
	}

	if err!=nil{
		log.Error(err)
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetOptUpload,err.Error())
		}
	}else{
		items:=strings.Split(url,",")
		urlInfo:=gabs.New()
		urlInfo.Set(items[0],"group")
		urlInfo.Set(items[1],"storage")
		urlInfo.Set(items[2],"volume")
		urlInfo.Set(items[3],"key")
		urlInfo.Set(items[4],"cookie")

		json:=gabs.New()
		json.Set(urlInfo.Data(),"urlInfo")

		url=base64.StdEncoding.EncodeToString([]byte(url))
		json.Set(url,"url")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
