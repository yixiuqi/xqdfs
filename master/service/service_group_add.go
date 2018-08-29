package service

import (
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	configuredef "xqdfs/configure/defines"

	"github.com/Jeffail/gabs"
)

func ServiceGroupAdd(context *Context,m map[string]interface{}) interface{}{
	groups,err:=context.ConfigureServer.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGetAll,err.Error())
	}

	id:=int32(0)
	for _,g:=range groups {
		if id<g.Id {
			id=g.Id
		}
	}

	newGroup:=configuredef.NewGroupDal()
	newGroup.Id=id+1
	newGroup.ReadOnly=true
	log.Debugf("group[%d] add",newGroup.Id)
	err=context.ConfigureServer.GroupAdd(newGroup)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupAdd,err.Error())
	}

	json:=gabs.New()
	json.Set(id,"newId")
	return helper.ResultBuildWithBody(constant.Success,json)
}
