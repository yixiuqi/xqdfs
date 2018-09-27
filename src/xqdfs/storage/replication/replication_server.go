package replication

import (
	"time"
	"sync"
	"fmt"

	"xqdfs/proxy"
	"xqdfs/configure"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/storage/conf"
	"xqdfs/storage/store"
	"xqdfs/storage/replication/process"
	configuredef "xqdfs/configure/defines"
)

type ReplicationServer struct {
	conf *conf.Config
	s *store.Store
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage

	wg sync.WaitGroup
	isRun bool
	signal chan int
	taskLock sync.RWMutex
	task map[int32]*process.ReplicationTask
}

func NewReplicationServer(conf *conf.Config,s *store.Store,configureServer *configure.ConfigureServer) (*ReplicationServer,error) {
	discoveryServer,err:= discovery.NewDiscoveryServer(configureServer)
	if err != nil {
		log.Errorf("create discovery server error[%v]",err)
		return nil,err
	}

	proxyStorage,err:= proxy.NewProxyStorage()
	if err != nil {
		log.Errorf("create proxy error[%v]",err)
		return nil,err
	}

	sync:=&ReplicationServer{
		conf:conf,
		s:s,
		configureServer:configureServer,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		isRun:true,
		signal:make(chan int, 1),
		task:make(map[int32]*process.ReplicationTask),
	}
	sync.probe()
	go sync.start()
	return sync,nil
}

func (this *ReplicationServer) start() {
	for this.isRun {
		func(){
			defer helper.HandleErr()
			this.probe()
			select {
			case <-time.After(3 * time.Second):
			case <-this.signal:
			}
		}()
	}
	this.wg.Done()
}

func (this *ReplicationServer) Stop() {
	log.Info("ReplicationServer stop")
	if this.isRun == false {
		return
	}

	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)

	this.taskLock.RLock()
	for _,v:=range this.task {
		v.Stop()
	}
	this.taskLock.RUnlock()

	if this.discoveryServer!=nil {
		this.discoveryServer.Stop()
	}

	if this.proxyStorage!=nil {
		this.proxyStorage.Stop()
	}
}

func (this *ReplicationServer) dump(){
	this.taskLock.Lock()
	defer this.taskLock.Unlock()
	log.Debug(this.task)
}

func (this *ReplicationServer) probe() {
	me:=fmt.Sprintf("%s:%d",this.conf.Server.Host,this.conf.Server.Port)
	var storages []*configuredef.StorageDal

	groups,err:=this.configureServer.GroupGetAll()
	if err!=nil{
		log.Error(err)
		return
	}

	for _,g:=range groups {
		for _,s:=range g.Storage {
			if s.Addr==me {
				storages=g.Storage
				break
			}
		}

		if storages!=nil{
			break
		}
	}

	if storages==nil {
		return
	}

	neighbor:=make([]*configuredef.StorageDal,0)
	for _,s:=range storages {
		if s.Addr!=me {
			neighbor=append(neighbor,s)
		}
	}

	this.taskLock.Lock()
	defer this.taskLock.Unlock()
	this.clear(neighbor)
	this.add(neighbor)
	this.updateState()
}

func (this *ReplicationServer) add(neighbor []*configuredef.StorageDal) {
	for _,n:=range neighbor {
		_,ok:=this.task[n.Id]
		if ok==true{
			continue
		}

		task,err:=process.NewReplicationTask(this.conf.Replication.Path,this.s,this.proxyStorage,n.Id,n.Addr)
		if err!=nil{
			log.Error(err)
		}else{
			this.task[n.Id]=task
		}
	}
}

func (this *ReplicationServer) clear(neighbor []*configuredef.StorageDal) {
	isExit:=true
	for isExit {
		isExit=false

		for k,v:=range this.task {
			isFind:=false
			for _,n:=range neighbor {
				if n.Id==k {
					isFind=true
					break
				}
			}

			if isFind==false{
				v.Destroy()
				delete(this.task,k)
				isExit=true
				break
			}
		}
	}
}

func (this *ReplicationServer) updateState() {
	storages:=this.discoveryServer.Storages()
	for _,v:=range this.task {
		for _,s:=range storages {
			if v.StorageId==s.Id{
				v.Online=s.Online
				break
			}
		}
	}
}

func (this *ReplicationServer) Replication(p process.Replication) bool {
	if this.isRun == false {
		log.Debug("this.isRun == false")
		return false
	}

	task:=make(map[int32]*process.ReplicationTask)
	this.taskLock.RLock()
	for k,v:=range this.task {
		task[k]=v
	}
	this.taskLock.RUnlock()
	if len(task)>0 {
		p.Process(task)
	}

	return true
}
