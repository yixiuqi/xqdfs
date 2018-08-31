package clear

import (
	"sync"
	"time"
	"math"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/utils/helper"
	"xqdfs/storage/block"
	"xqdfs/master/resource/usage"
)

const(
	OrderClearThreshold = "OrderClearThreshold"		// 最少预留多少卷 default:5
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
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	wg sync.WaitGroup
	isRun bool
	signal chan int
	orderClearThreshold int
}

func NewClearTimeOld(configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) (*ClearTimeOld,error) {
	orderClearThreshold:=5
	value,err:=configureServer.ParamGet(OrderClearThreshold)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=configureServer.ParamSet(OrderClearThreshold,"5")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		orderClearThreshold,err=helper.StringToInt(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	log.Infof("%s[%d]",OrderClearThreshold,orderClearThreshold)
	t:=&ClearTimeOld{
		configureServer:configureServer,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		signal:make(chan int, 1),
		isRun:true,
		orderClearThreshold:orderClearThreshold,
	}
	go t.task()
	return t,nil
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
	if free>this.orderClearThreshold {
		return
	}

	log.Debugf("auto clear group[%d] storage[%d][%s] volume[%d][%s]",groupId,storageId,storageAddr,volumeId,timeString)
	err:=this.proxyStorage.StorageVolumeClear(storageAddr,volumeId,true)
	if err!=nil {
		log.Errorf("send auto clear command error[%v]",err)
	}
}

