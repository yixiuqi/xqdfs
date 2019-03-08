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
	plugin.PluginAddService(constant.CmdGroupReadOnly,ServiceGroupReadOnly)
}

type RequestGroupReadOnly struct {
	GroupId int32	`json:"id"`
	ReadOnly bool	`json:"readOnly"`
}
func ServiceGroupReadOnly(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestGroupReadOnly{}
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
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	group.ReadOnly=req.ReadOnly
	err=conf.GroupEdit(group)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupEdit,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
