package service

import (
	"math"

	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/discovery"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageGetAll,ServiceStorageGetAll)
}

func ServiceStorageGetAll(m map[string]interface{}) interface{}{
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		log.Errorf("%s no support",plugin.PluginDiscoveryServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

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

	storagesDal,err:=conf.StorageGetAll()
	if err!=nil{
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageGetAll,err.Error())
	}
	storages:=discoveryServer.Storages()

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

		su:=usage.GetStorageUsageFromArray(storages,storagesDal[pos].Id)
		if su==nil{
			jsonStorage.Set("","total")
			jsonStorage.Set("","used")
			jsonStorage.Set("","util")
			jsonStorage.Set("","imageCount")
			jsonStorage.Set("","imageDelCount")
		}else{
			vTotal:=float32(math.Trunc(float64(su.Total)/1024/1024/1024*1e3) * 1e-3)
			jsonStorage.Set(vTotal,"total")
			vUsed:=float32(math.Trunc(float64(su.Used)/1024/1024/1024*1e3) * 1e-3)
			jsonStorage.Set(vUsed,"used")
			jsonStorage.Set(su.Util,"util")
			jsonStorage.Set(su.ImageCount,"imageCount")
			jsonStorage.Set(su.ImageDelCount,"imageDelCount")
		}
		jsonStorages.ArrayAppend(jsonStorage.Data(),"rows")
	}
	return helper.ResultBuildWithBody(constant.Success,jsonStorages)
}
