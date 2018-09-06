package service

import (
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageRemove,ServiceStorageRemove)
}

func ServiceStorageRemove(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var id int32
	value,ok:=m["id"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			id=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"id missing")
	}

	err:=conf.StorageRemove(id)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageRemove,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
