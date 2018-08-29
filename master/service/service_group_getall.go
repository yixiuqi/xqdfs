package service

import (
	"sort"
	"math"
	"fmt"

	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/master/resource/usage"
	configuredef "xqdfs/configure/defines"

	"github.com/Jeffail/gabs"
)

func ServiceGroupGetAll(context *Context,m map[string]interface{}) interface{}{
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

	groups,err:=context.ConfigureServer.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetGroupGetAll,err.Error())
	}

	//sort
	for _,g:=range groups {
		sort.Sort(configuredef.StorageDalSortById(g.Storage))
	}
	sort.Sort(configuredef.GroupDalSortById(groups))
	groupsUsage:=usage.GetGroupsUsage(context.DiscoveryServer.Groups())

	jsonGroups:=gabs.New()
	jsonGroups.Array("rows")
	jsonGroups.Set(len(groups),"total")
	start:=(page-1)*rows
	end:=start+rows

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
			info+="写吞吐:0 MB/s"+"<br>"
			info+="读吞吐:0 MB/s"
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
					info+=fmt.Sprintf("写吞吐:%v MB/s<br>",vWriteFlow)
					vReadFlow:=float32(math.Trunc(float64(su.ReadFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("读吞吐:%v MB/s<br>",vReadFlow)
					jsonGroup.Set(info,"storageThroughput")
				}else{
					jsonGroup.Set("","storageOnline")
					info:="总容量:0 GB"+"<br>"
					info+="已使用:0 GB"+"<br>"
					info+="使用率:0"
					jsonGroup.Set(info,"storageInfo")
					info="写TPS:0"+"<br>"
					info+="读QPS:0"+"<br>"
					info+="写吞吐:0 MB/s"+"<br>"
					info+="读吞吐:0 MB/s"
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
					info+=fmt.Sprintf("写吞吐:%v MB/s<br>",vWriteFlow)
					vReadFlow:=float32(math.Trunc(float64(gu.ReadFlow)/1024/1024*1e3) * 1e-3)
					info+=fmt.Sprintf("读吞吐:%v MB/s<br>",vReadFlow)
					jsonGroup.Set(info,"groupThroughput")
				}else{
					info:="总容量:0 GB"+"<br>"
					info+="已使用:0 GB"+"<br>"
					info+="使用率:0"
					jsonGroup.Set(info,"groupInfo")
					info="写TPS:0"+"<br>"
					info+="读QPS:0"+"<br>"
					info+="写吞吐:0 MB/s"+"<br>"
					info+="读吞吐:0 MB/s"
					jsonGroup.Set(info,"groupThroughput")
				}

				jsonGroups.ArrayAppend(jsonGroup.Data(),"rows")
			}
		}
	}
	return helper.ResultBuildWithBody(constant.Success,jsonGroups)
}
