package service

import (
	"strings"
	"encoding/base64"

	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/master/strategy"
	"xqdfs/utils/plugin"

	"github.com/Jeffail/gabs"
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
func ServiceOptUpload(m map[string]interface{}) interface{}{
	var strategyServer *strategy.AllocStrategyServer
	if s:=plugin.PluginGetObject(plugin.PluginStrategyServer);s==nil {
		log.Errorf("%s no support",plugin.PluginStrategyServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		strategyServer=s.(*strategy.AllocStrategyServer)
	}

	var img []byte
	value,ok:=m["img"]
	if ok {
		img=helper.ImageGet("",value.(string))
		if img==nil||len(img)==0 {
			return helper.ResultBuildWithExtInfo(errors.RetImageData,errors.ErrImageData.Error())
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"img missing")
	}

	key:=helper.KeyGenerate()
	url,err:=strategyServer.Write(key,constant.Cookie,img)
	for err==errors.ErrNeedleExist {
		log.Error(err," try again")
		key=helper.KeyGenerate()
		url,err=strategyServer.Write(key,constant.Cookie,img)
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
		urlInfo.Set(helper.TimeFromKey(key),"time")

		json:=gabs.New()
		json.Set(urlInfo.Data(),"urlInfo")

		url=base64.StdEncoding.EncodeToString([]byte(url))
		json.Set(url,"url")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}