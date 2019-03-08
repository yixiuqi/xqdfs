package usage

import (
	"math"
	"sort"
	"bytes"
	"encoding/json"

	"xqdfs/utils/log"
	"xqdfs/discovery/defines"

	"github.com/Jeffail/gabs"
)

type GroupsUsage struct {
	Usage []*GroupUsage	`json:"usage_groups"`
	Total int64			`json:"total"`	//总空间
	Used int64			`json:"used"`		//已经使用空间
	Util float32		`json:"util"`		//使用率
}

type GroupUsage struct {
 	Id int32						`json:"id"`
	ReadOnly bool					`json:"read_only"`
	Total int64						`json:"total"`	//总空间
 	Used int64						`json:"used"`		//已经使用空间
	Util float32					`json:"util"`		//使用率
	ReadQps uint64					`json:"readQPS"`
	ReadFlow uint64					`json:"readFlow"`
	WriteTps uint64					`json:"writeTPS"`
	WriteFlow uint64				`json:"writeFlow"`
	StorageUsage []*StorageUsage	`json:"storage"`
	ImageCount uint64				`json:"imageCount"`
	ImageDelCount uint64 			`json:"imageDelCount"`
}

/*

 */
func (this *GroupsUsage) GetGroupUsage(groupId int32) *GroupUsage {
	for _,g:=range this.Usage {
		if g.Id==groupId {
			return g
		}
	}

	return nil
}

func (this *GroupsUsage) GetStorageUsage(groupId int32,storageId int32) *StorageUsage {
	for _,g:=range this.Usage {
		if g.Id==groupId {
			for _,s:=range g.StorageUsage {
				if s.Id==storageId {
					return s
				}
			}
		}
	}

	return nil
}

func (this *GroupsUsage) Json() (*gabs.Container,error) {
	j,err:=json.Marshal(this)
	if err!=nil{
		return nil,err
	}

	dec := json.NewDecoder(bytes.NewBuffer(j))
	dec.UseNumber()
	json,err:=gabs.ParseJSONDecoder(dec)
	if err!=nil{
		return nil,err
	}else{
		return json,nil
	}
}

//获取所有组使用情况
func GetGroupsUsage(groups []*defines.Group) *GroupsUsage {
	if len(groups)==0 {
		log.Debug("groups is null")
		return nil
	}

	groupsUsage:=&GroupsUsage{
		Usage:make([]*GroupUsage,0),
	}

	for _,g:=range groups {
		gu:=&GroupUsage{
			Id:g.Id,
			ReadOnly:g.ReadOnly,
			StorageUsage:make([]*StorageUsage,0),
		}
		for _,s:=range g.Storage {
			su:=GetStorageUsage(s)
			if su!=nil{
				gu.Total+=su.Total
				gu.Used+=su.Used
				gu.ReadQps+=su.ReadQps
				gu.WriteTps+=su.WriteTps
				gu.ReadFlow+=su.ReadFlow
				gu.WriteFlow+=su.WriteFlow
				gu.ImageCount+=su.ImageCount
				gu.ImageDelCount+=su.ImageDelCount
				gu.StorageUsage=append(gu.StorageUsage,su)
			}
		}

		if len(g.Storage)!=0 {
			gu.Total=gu.Total/int64(len(g.Storage))
			gu.Used=gu.Used/int64(len(g.Storage))
			gu.ImageCount=gu.ImageCount/uint64(len(g.Storage))
			gu.ImageDelCount=gu.ImageDelCount/uint64(len(g.Storage))
		}
		if gu.Total==0{
			gu.Util=0
		}else{
			v:=float64(gu.Used)/float64(gu.Total)
			gu.Util=float32(math.Trunc(v*1e8+0.5) * 1e-8)
		}

		groupsUsage.Total+=gu.Total
		groupsUsage.Used+=gu.Used
		groupsUsage.Usage=append(groupsUsage.Usage,gu)
	}

	if groupsUsage.Total==0{
		groupsUsage.Util=0
	}else{
		v:=float64(groupsUsage.Used)/float64(groupsUsage.Total)
		groupsUsage.Util=float32(math.Trunc(v*1e8+0.5) * 1e-8)
	}
	return groupsUsage
}

func GetGroupsUsageSortById(groups []*defines.Group) *GroupsUsage {
	groupUsage:=GetGroupsUsage(groups)
	if groupUsage==nil {
		return nil
	}

	//按id排序
	for _,g:=range groupUsage.Usage {
		for _,s:=range g.StorageUsage {
			//sort volume
			sort.Sort(VolumeUsageSortById(s.VolumeUsage))
		}
		//sort storage
		sort.Sort(StorageUsageSortById(g.StorageUsage))
	}
	//sort group
	sort.Sort(GroupUsageSortById(groupUsage.Usage))
	return groupUsage
}

func GetGroupsUsageSortBySize(groups []*defines.Group) *GroupsUsage {
	groupUsage:=GetGroupsUsage(groups)
	if groupUsage==nil {
		return nil
	}

	for _,g:=range groupUsage.Usage {
		for _,s:=range g.StorageUsage {
			sort.Sort(VolumeUsageSortBySize(s.VolumeUsage))
		}
		sort.Sort(StorageUsageSortBySize(g.StorageUsage))
	}
	sort.Sort(GroupUsageSortBySize(groupUsage.Usage))
	return groupUsage
}

func GetGroupsUsageSortByWriteTps(groups []*defines.Group) *GroupsUsage {
	groupUsage:=GetGroupsUsage(groups)
	if groupUsage==nil {
		return nil
	}

	for _,g:=range groupUsage.Usage {
		for _,s:=range g.StorageUsage {
			sort.Sort(VolumeUsageSortByWriteTps(s.VolumeUsage))
		}
		sort.Sort(StorageUsageSortByWriteTps(g.StorageUsage))
	}
	sort.Sort(GroupUsageSortByWriteTps(groupUsage.Usage))
	return groupUsage
}