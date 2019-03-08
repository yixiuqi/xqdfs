package service

import (
	"context"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/storage/store"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeAdd,ServiceVolumeAdd)
}

type RequestVolumeAdd struct {
	Vid int32 `json:"vid"`
}
func ServiceVolumeAdd(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeAdd{}
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}
	_,err=storage.AddVolume(req.Vid)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetVolumeAdd,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
