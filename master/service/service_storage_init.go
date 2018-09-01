package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/proxy"
	"xqdfs/utils/log"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageInit,ServiceStorageInit)
}

func ServiceStorageInit(m map[string]interface{}) interface{}{
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

	var storageId int32
	value,ok:=m["id"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			storageId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"id missing")
	}

	storage,err:=conf.StorageGet(storageId)
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
