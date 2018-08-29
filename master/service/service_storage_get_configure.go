package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
)

func ServiceStorageGetConfigure(context *Context,m map[string]interface{}) interface{}{
	var storageId int32
	value,ok:=m["storageId"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			storageId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"storageId missing")
	}

	storage,err:=context.ConfigureServer.StorageGet(storageId)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageGet,err.Error())
	}
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	json,err:=context.ProxyStorage.StorageGetConfigure(storage.Addr)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStoreConfigure,err.Error())
	}else{
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
