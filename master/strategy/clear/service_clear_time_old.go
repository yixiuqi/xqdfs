package clear

import (
	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"

	"github.com/Jeffail/gabs"
)

const(
	CmdClearTimeOldConfigGet 	= "/strategy/clear/timeold/config/get"
	CmdClearTimeOldConfigSet 	= "/strategy/clear/timeold/config/set"
)

var(
	clearTimeOld *ClearTimeOld
)

func ServiceClearTimeOldSetup(timeold *ClearTimeOld) {
	clearTimeOld=timeold
	plugin.PluginAddService(CmdClearTimeOldConfigGet,ServiceClearTimeOldConfigGet)
	plugin.PluginAddService(CmdClearTimeOldConfigSet,ServiceClearTimeOldConfigSet)
}

func ServiceClearTimeOldConfigGet(m map[string]interface{}) interface{}{
	json:=gabs.New()
	json.Set(clearTimeOld.ClearTimeOldClearThresholdGet(),"clearThreshold")
	json.Set(clearTimeOld.curAvailableVolume,"curAvailableVolume")
	json.Set(clearTimeOld.oldGroupId,"oldestGroupId")
	json.Set(clearTimeOld.oldStorageId,"oldestStorageId")
	json.Set(clearTimeOld.oldVolumeId,"oldestVolumeId")
	json.Set(clearTimeOld.oldTime,"oldestTime")
	return helper.ResultBuildWithBody(constant.Success,json)
}

func ServiceClearTimeOldConfigSet(m map[string]interface{}) interface{}{
	var clearThreshold int
	value,ok:=m["clearThreshold"]
	if ok {
		tmp,err:=helper.GetInt(value)
		if err==nil{
			clearThreshold=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"consumeCount missing")
	}

	log.Debugf("clearThreshold[%d]",clearThreshold)

	err:=clearTimeOld.ClearTimeOldClearThresholdSet(clearThreshold)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
