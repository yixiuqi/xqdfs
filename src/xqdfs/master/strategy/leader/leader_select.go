package leader

import (
	"sync"
	"time"
	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/configure"
	"xqdfs/utils/plugin"
	"xqdfs/errors"
)

const(
	Leader = "leader"
	LeaderTime = 10
)

type LeaderSelect struct {
	configureServer *configure.ConfigureServer
	meId string
	leaderId string

	wg sync.WaitGroup
	isRun bool
	signal chan int
}

func NewLeaderSelect() (*LeaderSelect,error) {
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return nil,errors.ErrNoSupport
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	l:=&LeaderSelect{
		configureServer:conf,
		signal:make(chan int, 1),
		isRun:true,
		meId:helper.UUIDBuild(),
	}
	ServiceLeaderSelectSetup(l)
	go l.task()
	return l,nil
}

func (this *LeaderSelect) task() {
	for this.isRun {
		this.process()
		select {
		case <-time.After(LeaderTime * time.Second / 3):
		case <-this.signal:
		}
	}

	this.wg.Done()
}

func (this *LeaderSelect) process() {
	defer helper.HandleErr()

	this.leaderId=""

	v,err:=this.configureServer.ParamGet(Leader)
	if err==nil{
		this.leaderId=v
		if this.meId==v{
			this.configureServer.ParamSetx(Leader,this.meId,LeaderTime)
		}
	}else{
		this.configureServer.ParamSetx(Leader,this.meId,LeaderTime)
		v,err:=this.configureServer.ParamGet(Leader)
		if err==nil{
			this.leaderId=v
		}
	}
}

func (this *LeaderSelect) IsLeader() bool {
	return this.meId==this.leaderId
}

func (this *LeaderSelect) LeaderId() string {
	return this.leaderId
}

func (this *LeaderSelect) Stop() {
	log.Info("LeaderSelect stop")
	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)
}