package service

import (
	"sort"
	"math"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdStorageVolumeGetAll,ServiceStorageVolumeGetAll)
}

type RequestStorageVolumeGetAll struct {
	Page int32 			`json:"page"`
	Rows int32 			`json:"rows"`
	StorageId int32 	`json:"storageId"`
	Sort string 		`json:"sort"`
}
func ServiceStorageVolumeGetAll(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestStorageVolumeGetAll{}
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

		if j.Exists("storageId") {
			v,err:=helper.StringToInt32(j.Path("storageId").Data().(string))
			if err==nil {
				req.StorageId=v
			}
		}

		if j.Exists("sort") {
			req.Sort=j.Path("sort").Data().(string)
		}
	}

	if req.Rows==0 {
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"missing parameter rows")
	}

	if req.Page==0 {
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"missing parameter page")
	}

	if req.StorageId==0 {
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"missing parameter storageId")
	}

	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		log.Errorf("%s no support",plugin.PluginDiscoveryServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

	su:=usage.GetStorageUsageFromArray(discoveryServer.Storages(),req.StorageId)
	if su==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	if req.Sort=="byId"{
		sort.Sort(usage.VolumeUsageSortById(su.VolumeUsage))
	}else if req.Sort=="byUtil"{
		sort.Sort(usage.VolumeUsageSortByUtil(su.VolumeUsage))
	}else{
		sort.Sort(usage.VolumeUsageSortById(su.VolumeUsage))
	}

	jsonAll:=gabs.New()
	jsonAll.Array("rows")
	jsonAll.Set(len(su.VolumeUsage),"total")
	start:=(req.Page-1)*req.Rows
	end:=start+req.Rows
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
		jsonV.Set(su.VolumeUsage[pos].LastTime,"lastTime")
		jsonAll.ArrayAppend(jsonV.Data(),"rows")
	}
	return helper.ResultBuildWithBody(constant.Success,jsonAll)
}
