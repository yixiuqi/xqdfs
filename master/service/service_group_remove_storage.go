package service

import (
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	configuredef "xqdfs/configure/defines"
)

func ServiceGroupRemoveStorage(context *Context,m map[string]interface{}) interface{}{
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

	group,err:=context.ConfigureServer.GroupGet(groupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	newStorageArray:=make([]*configuredef.StorageDal,0)
	for _,s:=range group.Storage {
		if s.Id!=storageId {
			newStorageArray=append(newStorageArray,s)
		}
	}

	group.Storage=newStorageArray
	err=context.ConfigureServer.GroupEdit(group)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupRemoveStorage,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
