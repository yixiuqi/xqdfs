package service

import (
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/configure/defines"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.HttpStorageAdd,ServiceStorageAdd)
}

func ServiceStorageAdd(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var addr string
	var desc string
	value,ok:=m["addr"]
	if ok {
		addr=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"addr missing")
	}
	if helper.HostAddrCheck(addr) == false{
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,"addr param error")
	}

	value,ok=m["desc"]
	if ok {
		desc=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"desc missing")
	}

	storages,err:=conf.StorageGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageGetAll,err.Error())
	}

	id:=int32(0)
	for _,s:=range storages {
		if id<s.Id {
			id=s.Id
		}
	}

	newStorage:=defines.NewStorageDal()
	newStorage.Id=id+1
	newStorage.Addr=addr
	newStorage.Desc=desc
	err=conf.StorageAdd(newStorage)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageAdd,err.Error())
	}

	json:=gabs.New()
	json.Set(id,"newId")
	return helper.ResultBuildWithBody(constant.Success,json)
}
