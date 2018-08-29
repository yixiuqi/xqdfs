package service

import (
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	configuredef "xqdfs/configure/defines"

	"github.com/Jeffail/gabs"
)

func ServiceStorageAdd(context *Context,m map[string]interface{}) interface{}{
	var addr string
	var desc string

	value,ok:=m["addr"]
	if ok {
		addr=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"addr missing")
	}
	if helper.HostAddrCheck(addr) == false{
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,"addr param error")
	}

	value,ok=m["desc"]
	if ok {
		desc=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"desc missing")
	}

	storages,err:=context.ConfigureServer.StorageGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageGetAll,err.Error())
	}

	id:=int32(0)
	for _,s:=range storages {
		if id<s.Id {
			id=s.Id
		}
	}

	newStorage:=configuredef.NewStorageDal()
	newStorage.Id=id+1
	newStorage.Addr=addr
	newStorage.Desc=desc
	err=context.ConfigureServer.StorageAdd(newStorage)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageAdd,err.Error())
	}

	json:=gabs.New()
	json.Set(id,"newId")
	return helper.ResultBuildWithBody(constant.Success,json)
}
