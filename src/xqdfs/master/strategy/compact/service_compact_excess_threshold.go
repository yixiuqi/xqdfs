package compact

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
	CmdCompactExcessThresholdGet 	= "/strategy/compact/threshold/config/get"
	CmdCompactExcessThresholdSet 	= "/strategy/compact/threshold/config/set"
)

func setupService() {
	plugin.PluginAddService(CmdCompactExcessThresholdGet,ServiceCompactExcessThresholdGet)
	plugin.PluginAddService(CmdCompactExcessThresholdSet,ServiceCompactExcessThresholdSet)
}

func ServiceCompactExcessThresholdGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	json:=gabs.New()
	json.Set(compactExcessThresholdMinCountGet(),"minCount")
	json.Set(compactExcessThresholdValueGet(),"threshold")
	return helper.ResultBuildWithBody(constant.Success,json)
}

type RequestCompactExcessThresholdSet struct {
	MinCount int64		`json:"minCount"`
	Threshold float64	`json:"threshold"`
}
func ServiceCompactExcessThresholdSet(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestCompactExcessThresholdSet{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	err=compactExcessThresholdMinCountSet(req.MinCount)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	err=compactExcessThresholdValueSet(req.Threshold)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
