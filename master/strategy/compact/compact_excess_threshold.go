package compact

import (
	"fmt"
	"time"
	"sync"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/configure"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/resource/usage"
)

const(
	CompactExcessThresholdValue = "CompactExcessThresholdValue"
	CompactExcessThresholdMinCount = "CompactExcessThresholdMinCount"
)

type CompactExcessThreshold struct {
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	wg sync.WaitGroup
	isRun bool
	signal chan int
	excessThreshold float64
	minCount int64
}

func NewCompactExcessThreshold() (*CompactExcessThreshold,error) {
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

	excessThreshold:=0.3
	value,err:=conf.ParamGet(CompactExcessThresholdValue)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=conf.ParamSet(CompactExcessThresholdValue,"0.3")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		excessThreshold,err=helper.StringToFloat64(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	minCount:=int64(10000)
	value,err=conf.ParamGet(CompactExcessThresholdMinCount)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=conf.ParamSet(CompactExcessThresholdMinCount,"10000")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		minCount,err=helper.StringToInt64(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	c:=&CompactExcessThreshold{
		configureServer:conf,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		signal:make(chan int, 1),
		isRun:true,
		excessThreshold:excessThreshold,
		minCount:minCount,
	}
	ServiceCompactExcessThresholdSetup(c)
	go c.task()
	return c,nil
}

func (this *CompactExcessThreshold) task() {
	for this.isRun {
		this.process()
		select {
		case <-time.After(10 * time.Second):
		case <-this.signal:
		}
	}

	this.wg.Done()
}

func (this *CompactExcessThreshold) process() {
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

	for _,g:=range u.Usage {
		for _,s:=range g.StorageUsage {
			for _,v:=range s.VolumeUsage {
				if v.Compact || v.ImageCount<uint64(this.minCount) {
					continue
				}

				var util float64
				if v.ImageCount ==0 {
					util=0
				}else{
					util=float64(v.ImageDelCount)/float64(v.ImageCount)
					//log.Debugf("storage[%s] volume[%d] del[%v][%v]",s.Addr,v.Id,util,this.excessThreshold)
				}
				if util>this.excessThreshold {
					log.Debugf("auto compact [%s] volume[%d] util[%v]",s.Addr,v.Id,util)
					this.proxyStorage.StorageVolumeCompact(s.Addr,v.Id,true)
				}
			}
		}
	}
}

func (this *CompactExcessThreshold) Stop() {
	log.Info("CompactExcessThreshold stop")
	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)
}

func (this *CompactExcessThreshold) CompactExcessThresholdValueGet() float64 {
	return this.excessThreshold
}

func (this *CompactExcessThreshold) CompactExcessThresholdValueSet(excessThreshold float64) error {
	err:=this.configureServer.ParamSet(CompactExcessThresholdValue,fmt.Sprintf("%v",excessThreshold))
	if err!=nil{
		log.Error(err)
		return err
	}else{
		this.excessThreshold=excessThreshold
		return nil
	}
}

func (this *CompactExcessThreshold) CompactExcessThresholdMinCountGet() int64 {
	return this.minCount
}

func (this *CompactExcessThreshold) CompactExcessThresholdMinCountSet(minCount int64) error {
	err:=this.configureServer.ParamSet(CompactExcessThresholdMinCount,fmt.Sprintf("%d",minCount))
	if err!=nil{
		log.Error(err)
		return err
	}else{
		this.minCount=minCount
		return nil
	}
}