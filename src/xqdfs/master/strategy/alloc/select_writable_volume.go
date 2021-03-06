package alloc

import (
	"sync"
	"sort"

	"xqdfs/errors"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/master/strategy/tool"
	"xqdfs/master/resource/usage"
	strategydef "xqdfs/master/strategy/defines"
)

type SelectWritableVolume struct {
	discoveryServer *discovery.DiscoveryServer
	ptrTime int64
	groupUsageLock sync.Mutex
	groupUsage *usage.GroupsUsage
}

func NewSelectWritableVolume(discoveryServer *discovery.DiscoveryServer) *SelectWritableVolume {
	item:=&SelectWritableVolume{
		discoveryServer:discoveryServer,
	}

	return item
}

func (this *SelectWritableVolume) getGroupUsage() *usage.GroupsUsage {
	curTime:=helper.CurrentTime()
	if curTime-this.ptrTime>3000||this.groupUsage==nil{
		this.ptrTime=curTime
		this.groupUsage=usage.GetGroupsUsageSortById(this.discoveryServer.Groups())
	}
	if this.groupUsage==nil{
		log.Debug("groupUsage is null")
		return nil
	}

	//洗牌
	//for i := len(this.groupUsage.Usage) - 1; i > 0; i-- {
	//	num := helper.Uint32n(uint32(i + 1))
	//	this.groupUsage.Usage[i], this.groupUsage.Usage[num] = this.groupUsage.Usage[num], this.groupUsage.Usage[i]
	//}

	return this.groupUsage
}

func (this *SelectWritableVolume) SelectWritableVolume(fileSize int32,removeVolumes []*strategydef.WritableVolume) (*strategydef.WritableVolume,error) {
	this.groupUsageLock.Lock()
	defer this.groupUsageLock.Unlock()

	orderMinFreeSpace:=allocOrderMinFreeSpaceGet()
	orderConsumeCount:=allocOrderConsumeCountGet()

	availableVolumeAll:=make([]*strategydef.WritableVolume,0)
	u:=this.getGroupUsage()
	if u==nil {
		return nil,errors.ErrGroupIsEmpty
	}

	size:=tool.FileSizeCalc(fileSize)

	for _,g:=range u.Usage {
		if g.ReadOnly == true {
			continue
		}

		tolerant:=len(g.StorageUsage)
		su:=make([]*usage.StorageUsage,0)
		id:=make([]int,0)
		for sid:=0;sid<tolerant;sid++ {
			if g.StorageUsage[sid].Online == true {
				su=append(su,g.StorageUsage[sid])
				id=append(id,sid)
			}
		}

		if len(su)==0 {
			log.Debugf("group[%d] all storage offline",g.Id)
			continue
		}

		//这里是个技巧性的东西，同个组多个Storage作为互备的时候，选择volume是间隔选择的，方便互相同步数据
		i := helper.Uint32n(uint32(len(su)))
		objStorage:=su[i]
		for vid:=id[i];vid<len(objStorage.VolumeUsage);vid+=tolerant {
			//正在压缩
			if objStorage.VolumeUsage[vid].Compact {
				continue
			}

			//过滤
			isRemove:=false
			for _,r:=range removeVolumes {
				if r.GroupId==g.Id &&r.StorageId==objStorage.Id && r.VolumeId==objStorage.VolumeUsage[vid].Id{
					isRemove=true
					break
				}
			}
			if isRemove {
				continue
			}

			free:=objStorage.VolumeUsage[vid].Total - objStorage.VolumeUsage[vid].Used - orderMinFreeSpace
			if int64(size) < free {
				availableVolume:=&strategydef.WritableVolume{
					GroupId:g.Id,
					StorageId:objStorage.Id,
					StorageAddr:objStorage.Addr,
					VolumeId:objStorage.VolumeUsage[vid].Id,
					Free:free,
				}
				availableVolumeAll=append(availableVolumeAll,availableVolume)
			}
		}
	}

	if len(availableVolumeAll)==0 {
		log.Debug(errors.ErrSuperBlockNoSpace)
		return nil,errors.ErrSuperBlockNoSpace
	}else{
		sort.Stable(strategydef.WritableVolumeSort(availableVolumeAll))
		if len(availableVolumeAll) > orderConsumeCount {
			availableVolumeAll=availableVolumeAll[:orderConsumeCount]
		}

		i := helper.Uint32n(uint32(len(availableVolumeAll)))
		available:=availableVolumeAll[i]
		for _,g:=range u.Usage {
			if g.Id!=available.GroupId {
				continue
			}

			for _,s:=range g.StorageUsage {
				if s.Id!=available.StorageId {
					continue
				}

				for _,v:=range s.VolumeUsage {
					if v.Id==available.VolumeId {
						v.Used+=int64(size)
						return available,nil
					}
				}
			}
		}
		return available,nil
	}
}

