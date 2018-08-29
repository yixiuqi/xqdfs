package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
)

func ServiceStorageInit(context *Context,m map[string]interface{}) interface{}{
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

	storage,err:=context.ConfigureServer.StorageGet(storageId)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageGet,err.Error())
	}
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	err=context.ProxyStorage.StorageInit(storage.Addr,true)
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
