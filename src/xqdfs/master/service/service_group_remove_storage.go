package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/configure/defines"

	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdGroupRemoveStorage,ServiceGroupRemoveStorage)
}

type RequestGroupRemoveStorage struct {
	GroupId int32		`json:"groupId"`
	StorageId int32		`json:"storageId"`
}
func ServiceGroupRemoveStorage(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestGroupRemoveStorage{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	group,err:=conf.GroupGet(req.GroupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	newStorageArray:=make([]*defines.StorageDal,0)
	for _,s:=range group.Storage {
		if s.Id!=req.StorageId {
			newStorageArray=append(newStorageArray,s)
		}
	}

	group.Storage=newStorageArray
	err=conf.GroupEdit(group)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupRemoveStorage,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
