package alloc

import (
	"fmt"
	"sync"
	"time"
	"container/list"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
)

type DeleteCmd struct {
	host string
	vid int32
	key int64
}

func (this *DeleteCmd) String() string {
	return fmt.Sprintf("%s %d %d",this.host,this.vid,this.key)
}

type UploadErrorProcess struct {
	proxyStorage *proxy.ProxyStorage
	lock sync.RWMutex
	items map[string]*list.List
	wg sync.WaitGroup
	isRun bool
	signal chan int
}

func NewUploadErrorProcess(proxyStorage *proxy.ProxyStorage) *UploadErrorProcess {
	u:=&UploadErrorProcess{
		proxyStorage:proxyStorage,
		signal:make(chan int, 1),
		isRun:true,
		items:make(map[string]*list.List),
	}
	go u.task()
	return u
}

func (this *UploadErrorProcess) task() {
	for this.isRun {
		if this.process() == true {
			continue
		}

		select {
		case <-time.After(5 * time.Second):
		case <-this.signal:
		}
	}
	this.wg.Done()
}

func (this *UploadErrorProcess) process() bool {
	defer helper.HandleErr()
	this.lock.RLock()
	defer this.lock.RUnlock()

	if len(this.items) == 0 {
		return false
	}

	isFind:=false
	for k,v:=range this.items {
		for e := v.Front(); e != nil; e = e.Next() {
			item:=e.Value.(*DeleteCmd)
			err:=this.proxyStorage.Delete(item.host,item.vid,item.key,true)
			log.Debugf("delete process[%v][%v]",item,err)
			if err==errors.ErrRpc{
				break
			}else{
				v.Remove(e)
				isFind=true
			}
		}

		if v.Len() == 0 {
			delete(this.items,k)
		}
	}

	return isFind
}

func (this *UploadErrorProcess) Stop() {
	log.Info("UploadErrorProcess stop")
	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)
}

func (this *UploadErrorProcess) RollBack(host string,vid int32,key int64) {
	log.Debugf("RollBack [%s][%d][%d]",host,vid,key)
	this.lock.Lock()
	defer this.lock.Unlock()

	l,ok:=this.items[host]
	if ok == false {
		l=list.New()
		this.items[host]=l
	}

	dc:=&DeleteCmd{
		host:host,
		vid:vid,
		key:key,
	}
	l.PushBack(dc)
}
