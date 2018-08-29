package service

import (
	"math"

	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	"xqdfs/master/resource/usage"
	discoverydef "xqdfs/discovery/defines"

	"github.com/Jeffail/gabs"
)

func ServiceStorageGetAll(context *Context,m map[string]interface{}) interface{}{
	var page int32
	var rows int32

	value,ok:=m["page"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			page=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"page missing")
	}

	value,ok=m["rows"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			rows=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"rows missing")
	}

	storagesDal,err:=context.ConfigureServer.StorageGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageGetAll,err.Error())
	}
	storages:=context.DiscoveryServer.Storages()

	jsonStorages:=gabs.New()
	jsonStorages.Array("rows")
	jsonStorages.Set(len(storagesDal),"total")
	start:=(page-1)*rows
	end:=start+rows

	for pos:=start;pos<end&&pos<int32(len(storagesDal));pos++ {
		jsonStorage:=gabs.New()
		jsonStorage.Set(storagesDal[pos].Id,"id")
		jsonStorage.Set(storagesDal[pos].Addr,"addr")
		jsonStorage.Set(storagesDal[pos].Desc,"desc")

		var storage *discoverydef.Storage
		for _,s:=range storages {
			if s.Id==storagesDal[pos].Id {
				storage=s
				break
			}
		}

		if storage==nil{
			jsonStorage.Set("","total")
			jsonStorage.Set("","used")
			jsonStorage.Set("","util")
			jsonStorage.Set("","imageCount")
		}else{
			su:=usage.GetStorageUsage(storage)
			vTotal:=float32(math.Trunc(float64(su.Total)/1024/1024/1024*1e3) * 1e-3)
			jsonStorage.Set(vTotal,"total")
			vUsed:=float32(math.Trunc(float64(su.Used)/1024/1024/1024*1e3) * 1e-3)
			jsonStorage.Set(vUsed,"used")
			jsonStorage.Set(su.Util,"util")
			jsonStorage.Set(su.ImageCount,"imageCount")
		}
		jsonStorages.ArrayAppend(jsonStorage.Data(),"rows")
	}
	return helper.ResultBuildWithBody(constant.Success,jsonStorages)
}
