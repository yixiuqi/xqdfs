package compact

import (
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"

	"github.com/Jeffail/gabs"
)

const(
	CmdCompactExcessThresholdGet 	= "/strategy/compact/threshold/config/get"
	CmdCompactExcessThresholdSet 	= "/strategy/compact/threshold/config/set"
)

var(
	excessThreshold *CompactExcessThreshold
)

func ServiceCompactExcessThresholdSetup(excess *CompactExcessThreshold) {
	excessThreshold=excess
	plugin.PluginAddService(CmdCompactExcessThresholdGet,ServiceCompactExcessThresholdGet)
	plugin.PluginAddService(CmdCompactExcessThresholdSet,ServiceCompactExcessThresholdSet)
}

func ServiceCompactExcessThresholdGet(m map[string]interface{}) interface{}{
	json:=gabs.New()
	json.Set(excessThreshold.CompactExcessThresholdMinCountGet(),"minCount")
	json.Set(excessThreshold.CompactExcessThresholdValueGet(),"threshold")
	return helper.ResultBuildWithBody(constant.Success,json)
}

func ServiceCompactExcessThresholdSet(m map[string]interface{}) interface{}{
	var minCount int64
	var threshold float64
	value,ok:=m["minCount"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			minCount=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"minCount missing")
	}

	value,ok=m["threshold"]
	if ok {
		tmp,err:=helper.GetFloat64(value)
		if err==nil{
			threshold=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"threshold missing")
	}

	log.Debugf("minCount[%d] threshold[%v]",minCount,threshold)

	err:=excessThreshold.CompactExcessThresholdMinCountSet(minCount)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	err=excessThreshold.CompactExcessThresholdValueSet(threshold)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
