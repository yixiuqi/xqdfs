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
	plugin.PluginAddService(constant.CmdGroupRemove,ServiceGroupRemove)
}

type RequestGroupRemove struct {
	GroupId int32 `json:"id"`
}
func ServiceGroupRemove(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestGroupRemove{}
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

	group,err:=conf.GroupGet(req.GroupId)
	if err!=nil{
		log.Warnf("GroupGet[%d] error[%v]",req.GroupId,err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	err=conf.GroupRemove(req.GroupId)
	if err!=nil {
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupRemove,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
