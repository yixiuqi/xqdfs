package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/utils/plugin"
	"xqdfs/storage/conf"
	"xqdfs/utils/helper"
)

func init() {
	plugin.PluginAddService(constant.CmdStoreConf,ServiceStoreConf)
}

func ServiceStoreConf(ctx context.Context,inv *plugin.Invocation) interface{}{
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
