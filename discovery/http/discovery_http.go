package http

import (
	"time"
	"sync"

	"xqdfs/utils/helper"
	"xqdfs/utils/stat"
	"xqdfs/utils/log"
	"xqdfs/discovery/defines"
	"xqdfs/configure"
	configuredef "xqdfs/configure/defines"
)

const(
	HeartBeatSeconds = 3
)

type DiscoveryHttp struct {
	groups []*defines.Group
	storages []*defines.Storage
	configureServer *configure.ConfigureServer
	isRun bool
	wg sync.WaitGroup
	heartbeat int
}

func NewDiscoveryHttp(configureServer *configure.ConfigureServer) (gmgr *DiscoveryHttp, err error) {
	gmgr=&DiscoveryHttp{
		configureServer:configureServer,
		heartbeat:HeartBeatSeconds,
	}
	gmgr.start()
	return
}

func (this *DiscoveryHttp) Groups() []*defines.Group {
	return this.groups
}

func (this *DiscoveryHttp) Storages() []*defines.Storage {
	return this.storages
}

func (this *DiscoveryHttp) start() {
	this.isRun=true
	go func(){
		for this.isRun {
			g,s,err:=this.task()
			if err==nil{
				this.groups=g
				this.storages=s
			}

			for i:=0;i<this.heartbeat&&this.isRun;i++{
				time.Sleep(time.Second)
			}
		}

		this.wg.Done()
	}()
}

func (this *DiscoveryHttp) Stop() {
	this.wg.Add(1)
	this.isRun=false
	this.wg.Wait()
	log.Info("DiscoveryHttp stop")
}

func (this* DiscoveryHttp) task() (groups []*defines.Group,storages []*defines.Storage,err error){
	defer helper.HandleErr()

	storages=make([]*defines.Storage,0)
	var storageDals []*configuredef.StorageDal
	storageDals,err=this.configureServer.StorageGetAll()
	if err==nil&&storageDals!=nil{
		for _,s:=range storageDals {
			var storage *defines.Storage
			storage,err=this.probeStorage(s.Addr)
			if err==nil{
				storage.Id=s.Id
				storage.Addr=s.Addr
				storage.Online=true
				storages=append(storages,storage)
			}
		}
	}

	var groupDals []*configuredef.GroupDal
	groupDals,err=this.configureServer.GroupGetAll()
	if err!=nil || groupDals==nil {
		return
	}

	groups=make([]*defines.Group,0)
	for _,gdal:=range groupDals {
		gnew:=defines.NewGroup()
		gnew.Id=gdal.Id
		gnew.ReadOnly=gdal.ReadOnly
		groups=append(groups,gnew)

		if gdal.Storage==nil||len(gdal.Storage)==0{
			continue
		}

		gnew.Storage=make([]*defines.Storage,0)
		gnew.Stat=&stat.Stats{}

		for _,s:=range gdal.Storage {
			var stmp *defines.Storage
			for _,storage:=range storages {
				if s.Id==storage.Id {
					stmp=storage
					break
				}
			}

			if stmp==nil{
				stmp:=defines.NewStorage()
				stmp.Id=s.Id
				stmp.Addr=s.Addr
				stmp.Online=false
				gnew.Storage=append(gnew.Storage,stmp)
			}else{
				gnew.Storage=append(gnew.Storage,stmp)
				gnew.Stat.Merge(stmp.Stat)
			}
		}
	}
	return
}

func (this* DiscoveryHttp) probeStorage(addr string) (s *defines.Storage,err error) {
	for i:=0;i<3;i++{
		s,err=HeartBeat(addr)
		if err==nil{
			return
		}
	}

	return
}
