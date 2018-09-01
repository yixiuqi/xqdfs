package order

import (
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/constant"
	"xqdfs/errors"
	"xqdfs/utils/log"

	"github.com/Jeffail/gabs"
)

const(
	CmdAllocOrderConfigGet 	= "/strategy/alloc/order/config/get"
	CmdAllocOrderConfigSet 	= "/strategy/alloc/order/config/set"
)

var(
	allocOrder *AllocOrder
)

func ServiceAllocOrderSetup(order *AllocOrder) {
	allocOrder=order
	plugin.PluginAddService(CmdAllocOrderConfigGet,ServiceAllocOrderConfigGet)
	plugin.PluginAddService(CmdAllocOrderConfigSet,ServiceAllocOrderConfigSet)
}

func ServiceAllocOrderConfigGet(m map[string]interface{}) interface{}{
	json:=gabs.New()
	json.Set(allocOrder.AllocOrderConsumeCountGet(),"consumeCount")
	json.Set(allocOrder.AllocOrderMinFreeSpaceGet(),"minFreeSpace")
	return helper.ResultBuildWithBody(constant.Success,json)
}

func ServiceAllocOrderConfigSet(m map[string]interface{}) interface{}{
	var consumeCount int
	var minFreeSpace int64
	value,ok:=m["consumeCount"]
	if ok {
		tmp,err:=helper.GetInt(value)
		if err==nil{
			consumeCount=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"consumeCount missing")
	}

	value,ok=m["minFreeSpace"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			minFreeSpace=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"minFreeSpace missing")
	}
	log.Debugf("consumeCount[%d] minFreeSpace[%d]",consumeCount,minFreeSpace)

	err:=allocOrder.AllocOrderConsumeCountSet(consumeCount)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	err=allocOrder.AllocOrderMinFreeSpaceSet(minFreeSpace)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
