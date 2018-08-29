package order

import (
	"sync"

	"xqdfs/master/resource/usage"
	"xqdfs/errors"
	"xqdfs/master/strategy/tool"
	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/master/conf"
	strategydef "xqdfs/master/strategy/defines"
	"xqdfs/discovery"
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

	for i := len(this.groupUsage.Usage) - 1; i > 0; i-- {
		num := helper.Uint32n(uint32(i + 1))
		this.groupUsage.Usage[i], this.groupUsage.Usage[num] = this.groupUsage.Usage[num], this.groupUsage.Usage[i]
	}

	return this.groupUsage
}

func (this *SelectWritableVolume) SelectWritableVolume(conf *conf.Config,fileSize int32,removeVolumes []*strategydef.WritableVolume) (*strategydef.WritableVolume,error) {
	this.groupUsageLock.Lock()
	defer this.groupUsageLock.Unlock()

	availableVolumeAll:=make([]*strategydef.WritableVolume,0)
	u:=this.getGroupUsage()
	if u==nil {
		log.Debug("groups is null")
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

		i := helper.Uint32n(uint32(len(su)))
		objStorage:=su[i]
		for vid:=id[i];vid<len(objStorage.VolumeUsage);vid+=tolerant {
			//remove
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

			free:=objStorage.VolumeUsage[vid].Total - objStorage.VolumeUsage[vid].Used - conf.AllocStrategy.OrderMinFreeSpace
			if int64(size) < free {
				availableVolume:=&strategydef.WritableVolume{
					GroupId:g.Id,
					StorageId:objStorage.Id,
					StorageAddr:objStorage.Addr,
					VolumeId:objStorage.VolumeUsage[vid].Id,
				}
				availableVolumeAll=append(availableVolumeAll,availableVolume)
			}

			if len(availableVolumeAll) >= conf.AllocStrategy.OrderConsumeCount {
				break
			}
		}

		if len(availableVolumeAll) >= conf.AllocStrategy.OrderConsumeCount {
			break
		}
	}

	if len(availableVolumeAll)==0 {
		log.Debug(errors.ErrSuperBlockNoSpace)
		return nil,errors.ErrSuperBlockNoSpace
	}else{
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

