package strategy

import (
	"xqdfs/utils/log"
	"xqdfs/master/strategy/alloc"
	"xqdfs/master/strategy/clear"
	"xqdfs/master/strategy/leader"
	"xqdfs/master/strategy/defines"
	"xqdfs/master/strategy/compact"

	"github.com/Jeffail/gabs"
)

type AllocStrategyServer struct {
	leaderSelect defines.Leader
	allocStrategy defines.AllocStrategy
	clearStrategy defines.ClearStrategy
	compactStrategy defines.CompactStrategy
}

func NewAllocStrategyServer() (*AllocStrategyServer,error){
	leaderSelect,err:=leader.NewLeaderSelect()
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	alloc,err:=order.NewAllocOrder()
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	clear,err:=clear.NewClearTimeOld(leaderSelect)
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	comp,err:=compact.NewCompactExcessThreshold(leaderSelect)
	if err!=nil{
		log.Error(err)
		return nil,err
	}

	s:=&AllocStrategyServer{
		leaderSelect:leaderSelect,
		allocStrategy:alloc,
		clearStrategy:clear,
		compactStrategy:comp,
	}

	return s,nil
}

func (this *AllocStrategyServer) Write(key int64,cookie int32,body *gabs.Container) (string,error) {
	return this.allocStrategy.Write(key,cookie,body)
}

func (this *AllocStrategyServer) Read(url string) ([]byte,error) {
	return this.allocStrategy.Read(url)
}

func (this *AllocStrategyServer) Delete(url string) error {
	return this.allocStrategy.Delete(url)
}

func (this *AllocStrategyServer) Stop() {
	this.allocStrategy.Stop()
	this.clearStrategy.Stop()
	this.compactStrategy.Stop()
	this.leaderSelect.Stop()
}
