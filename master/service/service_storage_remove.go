package service

import (
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
)

func ServiceStorageRemove(context *Context,m map[string]interface{}) interface{}{
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

	err:=context.ConfigureServer.StorageRemove(id)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageRemove,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
