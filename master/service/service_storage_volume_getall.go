package service

import (
	"sort"
	"math"

	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
)

func ServiceStorageVolumeGetAll(context *Context,m map[string]interface{}) interface{}{
	var groupId int32
	var storageId int32
	var page int32
	var rows int32
	var sortType string

	value,ok:=m["groupId"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			groupId=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"groupId missing")
	}

	value,ok=m["storageId"]
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

	gu:=usage.GetGroupsUsage(context.DiscoveryServer.Groups())
	if gu==nil{
		return helper.ResultBuildWithExtInfo(errors.RetGroupIsEmpty,errors.ErrGroupIsEmpty.Error())
	}

	storage:=gu.GetStorageUsage(groupId,storageId)
	if storage==nil{
		return helper.ResultBuildWithExtInfo(errors.RetStorageNotExist,errors.ErrStorageNotExist.Error())
	}

	if sortType=="byId"{
		sort.Sort(usage.VolumeUsageSortById(storage.VolumeUsage))
	}else if sortType=="byUtil"{
		sort.Sort(usage.VolumeUsageSortByUtil(storage.VolumeUsage))
	}else{
		sort.Sort(usage.VolumeUsageSortById(storage.VolumeUsage))
	}

	jsonAll:=gabs.New()
	jsonAll.Array("rows")
	jsonAll.Set(len(storage.VolumeUsage),"total")
	start:=(page-1)*rows
	end:=start+rows
	for pos:=start;pos<end&&pos<int32(len(storage.VolumeUsage));pos++ {
		jsonV:=gabs.New()
		jsonV.Set(storage.VolumeUsage[pos].Id,"id")
		vTotal:=float32(math.Trunc(float64(storage.VolumeUsage[pos].Total)/1024/1024/1024*1e3) * 1e-3)
		jsonV.Set(vTotal,"total")
		vUsed:=float32(math.Trunc(float64(storage.VolumeUsage[pos].Used)/1024/1024/1024*1e3) * 1e-3)
		jsonV.Set(vUsed,"used")
		jsonV.Set(storage.VolumeUsage[pos].Util,"util")
		jsonV.Set(storage.VolumeUsage[pos].WriteTps,"writeTps")
		jsonV.Set(storage.VolumeUsage[pos].ReadQps,"readQps")
		jsonV.Set(storage.VolumeUsage[pos].Compact,"compact")
		jsonV.Set(storage.VolumeUsage[pos].ImageCount,"imageCount")
		jsonV.Set(storage.VolumeUsage[pos].ImageDelCount,"imageDelCount")
		jsonAll.ArrayAppend(jsonV.Data(),"rows")
	}
	return helper.ResultBuildWithBody(constant.Success,jsonAll)
}
