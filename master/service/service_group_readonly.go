package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
)

func ServiceGroupReadOnly(context *Context,m map[string]interface{}) interface{}{
	var groupId int32
	var readOnly bool

	value,ok:=m["id"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			groupId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"id missing")
	}

	value,ok=m["readOnly"]
	if ok {
		readOnly=(value=="true")
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"readOnly missing")
	}

	group,err:=context.ConfigureServer.GroupGet(groupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	group.ReadOnly=readOnly
	err=context.ConfigureServer.GroupEdit(group)
	if err!=nil {
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupEdit,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
