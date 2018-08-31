package service

import (
	"sort"
	"math"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/discovery"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.HttpStorageVolumeGetAll,ServiceStorageVolumeGetAll)
}

func ServiceStorageVolumeGetAll(m map[string]interface{}) interface{}{
	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

	var storageId int32
	var page int32
	var rows int32
	var sortType string
	value,ok:=m["storageId"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			storageId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"storageId missing")
	}

	value,ok=m["page"]
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

	value,ok=m["sort"]
	if ok {
		sortType=value.(string)
	}

	su:=usage.GetStorageUsageFromArray(discoveryServer.Storages(),storageId)
	if su==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	if sortType=="byId"{
		sort.Sort(usage.VolumeUsageSortById(su.VolumeUsage))
	}else if sortType=="byUtil"{
		sort.Sort(usage.VolumeUsageSortByUtil(su.VolumeUsage))
	}else{
		sort.Sort(usage.VolumeUsageSortById(su.VolumeUsage))
	}

	jsonAll:=gabs.New()
	jsonAll.Array("rows")
	jsonAll.Set(len(su.VolumeUsage),"total")
	start:=(page-1)*rows
	end:=start+rows
	for pos:=start;pos<end&&pos<int32(len(su.VolumeUsage));pos++ {
		jsonV:=gabs.New()
		jsonV.Set(su.VolumeUsage[pos].Id,"id")
		vTotal:=float32(math.Trunc(float64(su.VolumeUsage[pos].Total)/1024/1024/1024*1e3) * 1e-3)
		jsonV.Set(vTotal,"total")
		vUsed:=float32(math.Trunc(float64(su.VolumeUsage[pos].Used)/1024/1024/1024*1e3) * 1e-3)
		jsonV.Set(vUsed,"used")
		jsonV.Set(su.VolumeUsage[pos].Util,"util")
		jsonV.Set(su.VolumeUsage[pos].WriteTps,"writeTps")
		jsonV.Set(su.VolumeUsage[pos].ReadQps,"readQps")
		jsonV.Set(su.VolumeUsage[pos].Compact,"compact")
		jsonV.Set(su.VolumeUsage[pos].ImageCount,"imageCount")
		jsonV.Set(su.VolumeUsage[pos].ImageDelCount,"imageDelCount")
		jsonAll.ArrayAppend(jsonV.Data(),"rows")
	}
	return helper.ResultBuildWithBody(constant.Success,jsonAll)
}
