package alloc

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

const(
	CmdAllocOrderConfigGet 	= "/strategy/alloc/order/config/get"
	CmdAllocOrderConfigSet 	= "/strategy/alloc/order/config/set"
)

func setupService() {
	plugin.PluginAddService(CmdAllocOrderConfigGet,ServiceAllocOrderConfigGet)
	plugin.PluginAddService(CmdAllocOrderConfigSet,ServiceAllocOrderConfigSet)
}

func ServiceAllocOrderConfigGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	json:=gabs.New()
	json.Set(allocOrderConsumeCountGet(),"consumeCount")
	json.Set(allocOrderMinFreeSpaceGet(),"minFreeSpace")
	return helper.ResultBuildWithBody(constant.Success,json)
}

type RequestAllocOrderConfigSet struct {
	ConsumeCount int	`json:"consumeCount"`
	MinFreeSpace int64	`json:"minFreeSpace"`
}
func ServiceAllocOrderConfigSet(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestAllocOrderConfigSet{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	err=allocOrderConsumeCountSet(req.ConsumeCount)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	err=allocOrderMinFreeSpaceSet(req.MinFreeSpace)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
