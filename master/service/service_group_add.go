package service

import (
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/configure/defines"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.HttpGroupAdd,ServiceGroupAdd)
}

func ServiceGroupAdd(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	groups,err:=conf.GroupGetAll()
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

	newGroup:=defines.NewGroupDal()
	newGroup.Id=id+1
	newGroup.ReadOnly=true
	log.Debugf("group[%d] add",newGroup.Id)
	err=conf.GroupAdd(newGroup)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupAdd,err.Error())
	}

	json:=gabs.New()
	json.Set(id,"newId")
	return helper.ResultBuildWithBody(constant.Success,json)
}
