package service

import (
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
)

func ServiceGroupRemove(context *Context,m map[string]interface{}) interface{}{
	var groupId int32

	value,ok:=m["id"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			groupId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"id missing")
	}

	group,err:=context.ConfigureServer.GroupGet(groupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	err=context.ConfigureServer.GroupRemove(groupId)
	if err!=nil {
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupRemove,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
