package service

import (
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/configure/defines"
)

func init() {
	plugin.PluginAddService(constant.CmdGroupRemoveStorage,ServiceGroupRemoveStorage)
}

func ServiceGroupRemoveStorage(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var groupId int32
	var storageId int32
	value,ok:=m["groupId"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			groupId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"groupId missing")
	}

	value,ok=m["storageId"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			storageId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"storageId missing")
	}

	group,err:=conf.GroupGet(groupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	newStorageArray:=make([]*defines.StorageDal,0)
	for _,s:=range group.Storage {
		if s.Id!=storageId {
			newStorageArray=append(newStorageArray,s)
		}
	}

	group.Storage=newStorageArray
	err=conf.GroupEdit(group)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupRemoveStorage,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
