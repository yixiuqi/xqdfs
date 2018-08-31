package clear

import (
	"sync"
	"time"
	"math"

	"xqdfs/proxy"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/master/conf"
	"xqdfs/utils/helper"
	"xqdfs/storage/block"
	"xqdfs/master/resource/usage"
)

type VolumeItem struct {
	groupId int32
	storageId int32
	storageAddr string
	volumeId int32
	total int64
	used int64
}

type ClearTimeOld struct {
	conf *conf.Config
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	wg sync.WaitGroup
	isRun bool
	signal chan int
}

func NewClearTimeOld(conf *conf.Config,configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) *ClearTimeOld {
	t:=&ClearTimeOld{
		conf:conf,
		configureServer:configureServer,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		signal:make(chan int, 1),
		isRun:true,
	}
	go t.task()
	return t
}

func (this *ClearTimeOld) task() {
	for this.isRun {
		this.process()
		select {
		case <-time.After(10 * time.Second):
		case <-this.signal:
		}
	}

	this.wg.Done()
}

func (this *ClearTimeOld) Stop() {
	log.Info("ClearTimeOld stop")
	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)
}

func (this *ClearTimeOld) process() {
	defer helper.HandleErr()

	groups:=this.discoveryServer.Groups()
	if groups==nil {
		return
	}
	u:=usage.GetGroupsUsage(groups)
	if u==nil {
		return
	}
	if u.Total==0 {
		return
	}

	var groupId int32
	var storageId int32
	var storageAddr string
	var volumeId int32
	var timeMin int32 = math.MaxInt32
	var timeString string
	free:=0

	for _,g:=range u.Usage {
		if g.ReadOnly==true{
			continue
		}

		for _,s:=range g.StorageUsage {
			for _,v:=range s.VolumeUsage {
				if v.Used==block.HeaderSize {
					free++
				}

				time:=helper.TimeFromKey(v.LastKey)
				if time<timeMin && time !=0 {
					groupId=g.Id
					storageId=s.Id
					storageAddr=s.Addr
					volumeId=v.Id
					timeMin=time
					timeString=helper.TimeStringFromKey(v.LastKey)
				}
			}
		}
	}

	log.Debugf("available volume count[%d] util[%v]",free,u.Util)
	if free>this.conf.AllocStrategy.OrderClearThreshold {
		return
	}

	log.Debugf("auto clear group[%d] storage[%d][%s] volume[%d][%s]",groupId,storageId,storageAddr,volumeId,timeString)
	err:=this.proxyStorage.StorageVolumeClear(storageAddr,volumeId,true)
	if err!=nil {
		log.Errorf("send auto clear command error[%v]",err)
	}
}

