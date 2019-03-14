package compact

import (
	"time"
	"sync"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/configure"
	"xqdfs/master/resource/usage"
	"xqdfs/master/strategy/defines"
)

type CompactExcessThreshold struct {
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	leader defines.Leader
	wg sync.WaitGroup
	isRun bool
	signal chan int
}

func NewCompactExcessThreshold(leader defines.Leader) (*CompactExcessThreshold,error) {
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

	c:=&CompactExcessThreshold{
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		leader:leader,
		signal:make(chan int, 1),
		isRun:true,
	}

	setupConfigureSvr(conf)
	setupService()
	go c.task()
	return c,nil
}

func (this *CompactExcessThreshold) task() {
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
		if len(g.StorageUsage)==0 {
			continue
		}

		su:=g.StorageUsage[0]
		for _,v:=range su.VolumeUsage {
			if v.Compact || v.ImageCount<uint64(compactExcessThresholdMinCountGet()) || v.Util < CompactExcessThresholdMinUtil {
				continue
			}

			var util float64
			if v.ImageCount ==0 {
				util=0
			}else{
				util=float64(v.ImageDelCount)/float64(v.ImageCount)
				//log.Debugf("storage[%s] volume[%d] del[%v][%v]",s.Addr,v.Id,util,this.excessThreshold)
			}
			if util>compactExcessThresholdValueGet() {
				log.Debugf("auto compact [%s] volume[%d] util[%v]",su.Addr,v.Id,util)
				this.proxyStorage.StorageVolumeCompact(su.Addr,v.Id,true)
			}
		}
	}
}

func (this *CompactExcessThreshold) Stop() {
	log.Info("CompactExcessThreshold stop->")
	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)
	log.Info("CompactExcessThreshold stop-<")
}