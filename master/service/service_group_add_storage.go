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
	plugin.PluginAddService(constant.CmdGroupAddStorage,ServiceGroupAddStorage)
}

func ServiceGroupAddStorage(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var groupId int32
	var storageId int32
	var storageAddr string
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

	value,ok=m["storageAddr"]
	if ok {
		storageAddr=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"storageAddr missing")
	}
	if helper.HostAddrCheck(storageAddr) == false{
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,"storageAddr param error")
	}

	group,err:=conf.GroupGet(groupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	//判断节点是否已经使用
	groups,err:=conf.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGetAll,err.Error())
	}

	for _,g:=range groups {
		for _,s:=range g.Storage {
			if s.Id==storageId {
				return helper.ResultBuildWithExtInfo(errors.RetStorageExist,errors.ErrStorageExist.Error())
			}
		}
	}

	newStorage:=defines.NewStorageDal()
	newStorage.Id=storageId
	newStorage.Addr=storageAddr
	group.Storage=append(group.Storage,newStorage)
	err=conf.GroupEdit(group)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupAddStorage,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
