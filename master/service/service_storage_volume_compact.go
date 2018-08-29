package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
)

func ServiceStorageVolumeCompact(context *Context,m map[string]interface{}) interface{}{
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

	storage,err:=context.ConfigureServer.StorageGet(id)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageGet,err.Error())
	}
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	err=context.ProxyStorage.StorageVolumeCompact(storage.Addr,vid,true)
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
