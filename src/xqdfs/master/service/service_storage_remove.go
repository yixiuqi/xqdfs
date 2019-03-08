package service

import (
	"context"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageRemove,ServiceStorageRemove)
}

type RequestStorageRemove struct {
	Id int32 `json:"id"`
}
func ServiceStorageRemove(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestStorageRemove{}
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	err=conf.StorageRemove(req.Id)
	if err!=nil{
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageRemove,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
