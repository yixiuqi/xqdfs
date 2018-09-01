package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/proxy"
	"xqdfs/utils/plugin"
	"xqdfs/utils/log"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageVolumeCompact,ServiceStorageVolumeCompact)
}

func ServiceStorageVolumeCompact(m map[string]interface{}) interface{}{
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

	var id int32
	var vid int32
	value,ok:=m["id"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			id=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"id missing")
	}

	value,ok=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"vid missing")
	}

	storage,err:=conf.StorageGet(id)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageGet,err.Error())
	}
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	err=proxyStorage.StorageVolumeCompact(storage.Addr,vid,true)
	if err!=nil{
		e,ok:=err.(errors.Error)
		if ok {
			return helper.ResultBuildWithExtInfo(int32(e),err.Error())
		}else{
			return helper.ResultBuildWithExtInfo(errors.RetVolumeCompact,err.Error())
		}
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
