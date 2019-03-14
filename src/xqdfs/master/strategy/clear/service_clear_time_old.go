package clear

import (
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"

	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
)

const(
	CmdClearTimeOldConfigGet 	= "/strategy/clear/timeold/config/get"
	CmdClearTimeOldConfigSet 	= "/strategy/clear/timeold/config/set"
)

var(
	clearTimeOld *ClearTimeOld
)

//服务类初始化安装
//1.安装清除功能模块
//2.安装服务接口
func setupService(timeold *ClearTimeOld) {
	clearTimeOld=timeold
	plugin.PluginAddService(CmdClearTimeOldConfigGet,ServiceClearTimeOldConfigGet)
	plugin.PluginAddService(CmdClearTimeOldConfigSet,ServiceClearTimeOldConfigSet)
}

func ServiceClearTimeOldConfigGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	json:=gabs.New()
	json.Set(clearTimeOldClearThresholdGet(),"clearThreshold")
	json.Set(clearTimeOld.CurAvailableVolume,"curAvailableVolume")
	json.Set(clearTimeOld.OldGroupId,"oldestGroupId")
	json.Set(clearTimeOld.OldStorageId,"oldestStorageId")
	json.Set(clearTimeOld.OldVolumeId,"oldestVolumeId")
	json.Set(clearTimeOld.OldTime,"oldestTime")
	return helper.ResultBuildWithBody(constant.Success,json)
}

type RequestClearTimeOldConfigSet struct {
	ClearThreshold int `json:"clearThreshold"`
}
func ServiceClearTimeOldConfigSet(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestClearTimeOldConfigSet{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	err=clearTimeOldClearThresholdSet(req.ClearThreshold)
	if err!=nil{
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParamSet,err.Error())
	}
	return helper.ResultBuild(constant.Success)
}
