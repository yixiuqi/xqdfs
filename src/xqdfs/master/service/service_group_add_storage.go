package service

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/master/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/configure/defines"

	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdGroupAddStorage,ServiceGroupAddStorage)
}

type RequestGroupAddStorage struct {
	GroupId int32 			`json:"groupId"`
	StorageId int32 		`json:"storageId"`
	StorageAddr string 		`json:"storageAddr"`
}
func ServiceGroupAddStorage(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestGroupAddStorage{}
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

	if helper.HostAddrCheck(req.StorageAddr) == false{
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,"storageAddr param error")
	}

	group,err:=conf.GroupGet(req.GroupId)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGet,err.Error())
	}
	if group==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupNotExist,errors.ErrGroupNotExist.Error())
	}

	//判断节点是否已经使用
	groups,err:=conf.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGetAll,err.Error())
	}

	for _,g:=range groups {
		for _,s:=range g.Storage {
			if s.Id==req.StorageId {
				return helper.ResultBuildWithExtInfo(errors.RetStorageExist,errors.ErrStorageExist.Error())
			}
		}
	}

	newStorage:=defines.NewStorageDal()
	newStorage.Id=req.StorageId
	newStorage.Addr=req.StorageAddr
	group.Storage=append(group.Storage,newStorage)
	err=conf.GroupEdit(group)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupAddStorage,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
