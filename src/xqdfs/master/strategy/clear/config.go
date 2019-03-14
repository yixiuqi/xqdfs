package clear

import (
	"fmt"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/master/configure"
)

const(
	ClearTimeOldClearThreshold 	= "ClearTimeOldClearThreshold"		// 最少预留多少卷 default:5
	ClearTimeOldClearEnable 		= "ClearTimeOldClearEnable"			// 是否启动自动清除最老卷
)

var(
	configureSvr *configure.ConfigureServer
)

func setupConfigureSvr(conf *configure.ConfigureServer) {
	configureSvr=conf
}

func clearTimeOldClearThresholdGet() int {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	clearTimeOldThreshold:=5
	value,err:=configureSvr.ConfigGet(ClearTimeOldClearThreshold)
	if err!=nil {
		if err==errors.ErrParamNotExist {
			log.Debugf("clearTimeOldThreshold not exist,now set %d",clearTimeOldThreshold)
			err=configureSvr.ConfigSet(ClearTimeOldClearThreshold,fmt.Sprintf("%d",clearTimeOldThreshold))
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		v,err:=helper.StringToInt(value)
		if err!=nil{
			log.Warn(err)
		}else{
			clearTimeOldThreshold=v
		}
	}

	return clearTimeOldThreshold
}

func clearTimeOldClearThresholdSet(clearTimeOldThreshold int) error {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	err:=configureSvr.ConfigSet(ClearTimeOldClearThreshold,fmt.Sprintf("%d",clearTimeOldThreshold))
	if err!=nil{
		log.Warn(err)
		return err
	}else{
		return nil
	}
}

func clearTimeOldClearEnableGet() bool {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	clearTimeOldClearEnable:=false
	value,err:=configureSvr.ConfigGet(ClearTimeOldClearEnable)
	if err!=nil {
		if err==errors.ErrParamNotExist {
			err=configureSvr.ConfigSet(ClearTimeOldClearEnable,"false")
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		if value=="true" {
			clearTimeOldClearEnable=true
		}
	}

	return clearTimeOldClearEnable
}
