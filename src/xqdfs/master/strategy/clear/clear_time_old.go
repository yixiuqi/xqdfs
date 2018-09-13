package clear

import (
	"sync"
	"time"
	"math"
	"fmt"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/storage/block"
	"xqdfs/master/resource/usage"
	"xqdfs/master/strategy/defines"
)

const(
	ClearTimeOldClearThreshold = "ClearTimeOldClearThreshold"		// 最少预留多少卷 default:5
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
	leader defines.Leader
	wg sync.WaitGroup
	isRun bool
	signal chan int

	ClearTimeOldThreshold int
	CurAvailableVolume int
	OldGroupId int32
	OldStorageId int32
	OldStorageAddr string
	OldVolumeId int32
	OldTime string
}

func NewClearTimeOld(leader defines.Leader) (*ClearTimeOld,error) {
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return nil,errors.ErrNoSupport
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		log.Errorf("%s no support",plugin.PluginDiscoveryServer)
		return nil,errors.ErrNoSupport
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

	var proxyStorage *proxy.ProxyStorage
	if p:=plugin.PluginGetObject(plugin.PluginProxyStorage);p==nil {
		log.Errorf("%s no support",plugin.PluginProxyStorage)
		return nil,errors.ErrNoSupport
	}else{
		proxyStorage=p.(*proxy.ProxyStorage)
	}

	clearTimeOldThreshold:=5
	value,err:=conf.ParamGet(ClearTimeOldClearThreshold)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=conf.ParamSet(ClearTimeOldClearThreshold,"5")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		clearTimeOldThreshold,err=helper.StringToInt(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	log.Infof("%s[%d]",ClearTimeOldClearThreshold,clearTimeOldThreshold)
	t:=&ClearTimeOld{
		configureServer:conf,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		leader:leader,
		signal:make(chan int, 1),
		isRun:true,
		ClearTimeOldThreshold:clearTimeOldThreshold,
	}
	ServiceClearTimeOldSetup(t)
	go t.task()
	return t,nil
}

func (this *ClearTimeOld) task() {
	for this.isRun {
		if this.leader.IsLeader() {
			this.process()
		}else{
			log.Debug("not leader")
		}

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

	var timeMin int64 = math.MaxInt64
	freeVolume:=0
	availableGroup:=0
	for _,g:=range u.Usage {
		if g.ReadOnly==true{
			continue
		}
		availableGroup++

		for _,s:=range g.StorageUsage {
			for _,v:=range s.VolumeUsage {
				if v.Used==block.HeaderSize {
					freeVolume++
				}else{
					if v.LastWriteTime<timeMin && v.LastWriteTime !=0 {
						this.OldGroupId=g.Id
						this.OldStorageId=s.Id
						this.OldStorageAddr=s.Addr
						this.OldVolumeId=v.Id
						timeMin=v.LastWriteTime
						this.OldTime=v.LastTime
					}
				}
			}
		}
	}
	if availableGroup==0 {
		return
	}

	this.CurAvailableVolume=freeVolume
	if freeVolume>this.ClearTimeOldThreshold {
		return
	}

	log.Debugf("auto clear group[%d] storage[%d][%s] volume[%d][%s]",
		this.OldGroupId,this.OldStorageId,this.OldStorageAddr,this.OldVolumeId,this.OldTime)
	err:=this.proxyStorage.StorageVolumeClear(this.OldStorageAddr,this.OldVolumeId,true)
	if err!=nil {
		log.Errorf("send auto clear command error[%v]",err)
	}
}

func (this *ClearTimeOld) ClearTimeOldClearThresholdSet(clearTimeOldThreshold int) error {
	err:=this.configureServer.ParamSet(ClearTimeOldClearThreshold,fmt.Sprintf("%d",clearTimeOldThreshold))
	if err!=nil{
		log.Error(err)
		return err
	}else{
		this.ClearTimeOldThreshold=clearTimeOldThreshold
		return nil
	}
}

