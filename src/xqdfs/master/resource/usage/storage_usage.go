package usage

import (
	"math"

	"xqdfs/discovery/defines"
)

type StorageUsage struct {
	Id int32					`json:"id"`
	Addr string					`json:"addr"`
	Online bool					`json:"online"`
	Total int64					`json:"total"`	//总空间
	Used int64					`json:"used"`		//已经使用空间
	Util float32				`json:"util"`		//使用率
	ReadQps uint64				`json:"readQPS"`
	ReadFlow uint64				`json:"readFlow"`
	WriteTps uint64				`json:"writeTPS"`
	WriteFlow uint64			`json:"writeFlow"`
	VolumeUsage []*VolumeUsage	`json:"volume"`
	ImageCount uint64			`json:"imageCount"`
	ImageDelCount uint64 		`json:"imageDelCount"`
}

func GetStorageUsage(storage *defines.Storage) *StorageUsage {
	if storage==nil{
		return nil
	}

	su:= &StorageUsage{
		Id:storage.Id,
		Addr:storage.Addr,
		Online:storage.Online,
		VolumeUsage:make([]*VolumeUsage,0),
	}

	for _,v:=range storage.Volume {
		vu:=GetVolumeUsage(v)
		if vu!=nil{
			su.Total+=vu.Total
			su.Used+=vu.Used
			su.ImageCount+=vu.ImageCount
			su.ImageDelCount+=vu.ImageDelCount
			su.WriteTps+=vu.WriteTps
			su.ReadQps+=vu.ReadQps
			su.WriteFlow+=vu.WriteFlow
			su.ReadFlow+=vu.ReadFlow
			su.VolumeUsage=append(su.VolumeUsage,vu)
		}
	}

	if su.Total==0{
		su.Util=0
	}else{
		v:=float64(su.Used)/float64(su.Total)
		su.Util=float32(math.Trunc(v*1e6+0.5) * 1e-6)
	}

	return su
}

func GetStorageUsageFromArray(storages []*defines.Storage,id int32) *StorageUsage {
	var storage *defines.Storage
	for _,s:=range storages {
		if s.Id==id {
			storage=s
			break
		}
	}
	if storage==nil{
		return nil
	}else{
		return GetStorageUsage(storage)
	}
}
