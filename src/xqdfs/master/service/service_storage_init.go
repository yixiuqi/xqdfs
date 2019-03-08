package service

import (
	"context"
	"encoding/json"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageInit,ServiceStorageInit)
}

type RequestStorageInit struct {
	StorageId int32	`json:"id"`
}
func ServiceStorageInit(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestStorageInit{}
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

	var proxyStorage *proxy.ProxyStorage
	if p:=plugin.PluginGetObject(plugin.PluginProxyStorage);p==nil {
		log.Errorf("%s no support",plugin.PluginProxyStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		proxyStorage=p.(*proxy.ProxyStorage)
	}

	storage,err:=conf.StorageGet(req.StorageId)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageGet,err.Error())
	}
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	err=proxyStorage.StorageInit(storage.Addr,true)
	if err!=nil{
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetStoreInitFailed,err.Error())
		}
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
