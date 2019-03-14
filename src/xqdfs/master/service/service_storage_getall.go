package service

import (
	"math"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/master/configure"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageGetAll,ServiceStorageGetAll)
}

type RequestStorageGetAll struct {
	Page int32 			`json:"page"`
	Rows int32 			`json:"rows"`
}
func ServiceStorageGetAll(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestStorageGetAll{}
	if inv.ContentType==plugin.HttpTextPlain||inv.ContentType==plugin.HttpApplicationJson {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		err:=json.Unmarshal(inv.Body,req)
		if err!=nil {
			log.Warn(err)
			return helper.ResultBuild(errors.RetParameterError)
		}
	}else{
		j,err:=gabs.ParseJSON(inv.Body)
		if err!=nil {
			log.Warn(err)
			return helper.ResultBuild(errors.RetParameterError)
		}

		if j.Exists("page") {
			v,err:=helper.StringToInt32(j.Path("page").Data().(string))
			if err==nil {
				req.Page=v
			}
		}

		if j.Exists("rows") {
			v,err:=helper.StringToInt32(j.Path("rows").Data().(string))
			if err==nil {
				req.Rows=v
			}
		}
	}

	if req.Rows==0 {
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"missing parameter rows")
	}

	if req.Page==0 {
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"missing parameter page")
	}

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

	storagesDal,err:=conf.StorageGetAll()
	if err!=nil{
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetStorageGetAll,err.Error())
	}
	storages:=discoveryServer.Storages()

	jsonStorages:=gabs.New()
	jsonStorages.Array("rows")
	jsonStorages.Set(len(storagesDal),"total")
	start:=(req.Page-1)*req.Rows
	end:=start+req.Rows

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
