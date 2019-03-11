package service

import (
	"fmt"
	"sort"
	"math"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/configure/defines"
	"xqdfs/master/resource/usage"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdGroupGetAll,ServiceGroupGetAll)
}

type RequestGroupGetAll struct {
	Page int32 			`json:"page"`
	Rows int32 			`json:"rows"`
}
func ServiceGroupGetAll(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestGroupGetAll{}
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

	groups,err:=conf.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGetAll,err.Error())
	}

	//sort
	for _,g:=range groups {
		sort.Sort(defines.StorageDalSortById(g.Storage))
	}
	sort.Sort(defines.GroupDalSortById(groups))
	groupsUsage:=usage.GetGroupsUsage(discoveryServer.Groups())

	jsonGroups:=gabs.New()
	jsonGroups.Array("rows")
	jsonGroups.Set(len(groups),"total")
	start:=(req.Page-1)*req.Rows
	end:=start+req.Rows

	for pos:=start;pos<end&&pos<int32(len(groups));pos++ {
		if groups[pos].Storage==nil||len(groups[pos].Storage)==0{
			jsonGroup:=gabs.New()
			jsonGroup.Set(groups[pos].Id,"groupId")
			jsonGroup.Set(groups[pos].ReadOnly,"groupReadOnly")
			info:="总容量:0 GB"+"<br>"
			info+="已使用:0 GB"+"<br>"
			info+="使用率:0"
			jsonGroup.Set(info,"groupInfo")
			info="写TPS:0"+"<br>"
			info+="读QPS:0"+"<br>"
			info+="写入量:0 MB/s"+"<br>"
			info+="读取量:0 MB/s"
			jsonGroup.Set(info,"groupThroughput")
			jsonGroup.Set("","storageId")
			jsonGroup.Set("","storageAddr")
			jsonGroup.Set("","storageOnline")
			jsonGroup.Set("","storageInfo")
			jsonGroup.Set("","storageThroughput")
			jsonGroups.ArrayAppend(jsonGroup.Data(),"rows")
		}else{
			for _,s:=range groups[pos].Storage {
				jsonGroup:=gabs.New()
				jsonGroup.Set(groups[pos].Id,"groupId")
				jsonGroup.Set(groups[pos].ReadOnly,"groupReadOnly")
				jsonGroup.Set(s.Id,"storageId")
				jsonGroup.Set(s.Addr,"storageAddr")

				var su *usage.StorageUsage
				var gu *usage.GroupUsage

				if groupsUsage!=nil{
					su=groupsUsage.GetStorageUsage(groups[pos].Id,s.Id)
					gu=groupsUsage.GetGroupUsage(groups[pos].Id)
				}

				if su!=nil{
					jsonGroup.Set(su.Online,"storageOnline")
					vTotal:=float32(math.Trunc(float64(su.Total)/1024/1024/1024*1e3) * 1e-3)
					vUsed:=float32(math.Trunc(float64(su.Used)/1024/1024/1024*1e3) * 1e-3)
					info:=fmt.Sprintf("总容量:%v GB<br>",vTotal)
					info+=fmt.Sprintf("已使用:%v GB<br>",vUsed)
					info+=fmt.Sprintf("使用率:%v",su.Util)
					jsonGroup.Set(info,"storageInfo")
					info=fmt.Sprintf("写TPS:%v<br>",su.WriteTps)
					info+=fmt.Sprintf("读QPS:%v<br>",su.ReadQps)
					vWriteFlow:=float32(math.Trunc(float64(su.WriteFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("写入量:%v MB/s<br>",vWriteFlow)
					vReadFlow:=float32(math.Trunc(float64(su.ReadFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("读取量:%v MB/s<br>",vReadFlow)
					jsonGroup.Set(info,"storageThroughput")
				}else{
					jsonGroup.Set("","storageOnline")
					info:="总容量:0 GB"+"<br>"
					info+="已使用:0 GB"+"<br>"
					info+="使用率:0"
					jsonGroup.Set(info,"storageInfo")
					info="写TPS:0"+"<br>"
					info+="读QPS:0"+"<br>"
					info+="写入量:0 MB/s"+"<br>"
					info+="读取量:0 MB/s"
					jsonGroup.Set(info,"storageThroughput")
				}

				if gu!=nil{
					vTotal:=float32(math.Trunc(float64(gu.Total)/1024/1024/1024*1e3) * 1e-3)
					vUsed:=float32(math.Trunc(float64(gu.Used)/1024/1024/1024*1e3) * 1e-3)
					info:=fmt.Sprintf("总容量:%v GB<br>",vTotal)
					info+=fmt.Sprintf("已使用:%v GB<br>",vUsed)
					info+=fmt.Sprintf("使用率:%v",gu.Util)
					jsonGroup.Set(info,"groupInfo")
					info=fmt.Sprintf("写TPS:%v<br>",gu.WriteTps)
					info+=fmt.Sprintf("读QPS:%v<br>",gu.ReadQps)
					vWriteFlow:=float32(math.Trunc(float64(gu.WriteFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("写入量:%v MB/s<br>",vWriteFlow)
					vReadFlow:=float32(math.Trunc(float64(gu.ReadFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("读取量:%v MB/s<br>",vReadFlow)
					jsonGroup.Set(info,"groupThroughput")
				}else{
					info:="总容量:0 GB"+"<br>"
					info+="已使用:0 GB"+"<br>"
					info+="使用率:0"
					jsonGroup.Set(info,"groupInfo")
					info="写TPS:0"+"<br>"
					info+="读QPS:0"+"<br>"
					info+="写入量:0 MB/s"+"<br>"
					info+="读取量:0 MB/s"
					jsonGroup.Set(info,"groupThroughput")
				}

				jsonGroups.ArrayAppend(jsonGroup.Data(),"rows")
			}
		}
	}
	return helper.ResultBuildWithBody(constant.Success,jsonGroups)
}
