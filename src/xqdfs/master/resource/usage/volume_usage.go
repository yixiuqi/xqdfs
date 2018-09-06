package usage

import (
	"math"

	"xqdfs/discovery/defines"
	"xqdfs/storage/block"
	"xqdfs/utils/helper"
)

/*
 *	Volume使用情况描述
 */

type VolumeUsage struct {
	Id int32				`json:"id"`
	Total int64				`json:"total"`	//总空间
	Used int64				`json:"used"`		//已经使用空间
	Util float32			`json:"util"`		//使用率
	ReadQps uint64			`json:"readQPS"`
	ReadFlow uint64			`json:"readFlow"`
	WriteTps uint64			`json:"writeTPS"`
	WriteFlow uint64		`json:"writeFlow"`
	LastKey	int64			`json:"lastKey"`
	LastTime string			`json:"last_time"`
	Compact	bool			`json:"compact"`
	ImageCount uint64		`json:"imageCount"`
	ImageDelCount uint64 	`json:"imageDelCount"`
}

func GetVolumeUsage(volume *defines.Volume) *VolumeUsage {
	if volume==nil {
		return nil
	}

	v:= &VolumeUsage{
		Id:volume.Id,
	}

	v.Total=block.MaxSize
	v.Used=volume.Block.Size
	v.ReadQps=volume.Stat.ReadQPS
	v.ReadFlow=volume.Stat.ReadFlow
	v.WriteTps=volume.Stat.WriteTPS
	v.WriteFlow=volume.Stat.WriteFlow
	v.LastKey=volume.LastKey
	v.LastTime=helper.TimeStringFromKey(v.LastKey)
	v.Compact=volume.Compact
	v.ImageCount=volume.ImageCount
	v.ImageDelCount=volume.Stat.TotalDelProcessed

	if v.Total==0{
		v.Util=0
	}else{
		t:=float64(v.Used)/float64(v.Total)
		v.Util=float32(math.Trunc(t*1e6+0.5) * 1e-6)
	}

	return v
}
