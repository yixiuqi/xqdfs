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

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeAddFree,ServiceVolumeAddFree)
}

type RequestVolumeAddFree struct {
	Count int `json:"count"`
	Bdir string `json:"bdir"`
	Idir string `json:"idir"`
}
func ServiceVolumeAddFree(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeAddFree{Count:1}
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

	sn,err:=storage.AddFreeVolume(req.Count,req.Bdir,req.Idir)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetVolumeAddFree,err.Error())
	}else{
		json:=gabs.New()
		json.Set(sn,"count")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
