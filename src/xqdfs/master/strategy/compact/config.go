package compact

import (
	"fmt"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/master/configure"
)

const(
	CompactExcessThresholdValue = "CompactExcessThresholdValue"
	CompactExcessThresholdMinCount = "CompactExcessThresholdMinCount"
	CompactExcessThresholdMinUtil = 0.5
)

var(
	configureSvr *configure.ConfigureServer
)

func setupConfigureSvr(conf *configure.ConfigureServer) {
	configureSvr=conf
}

func compactExcessThresholdValueGet() float64 {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	excessThreshold:=0.3
	value,err:=configureSvr.ConfigGet(CompactExcessThresholdValue)
	if err!=nil{
		if err==errors.ErrParamNotExist {
			log.Debugf("compactExcessThresholdValue not exist,now set %v",excessThreshold)
			err=configureSvr.ConfigSet(CompactExcessThresholdValue,fmt.Sprintf("%v",excessThreshold))
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		v,err:=helper.StringToFloat64(value)
		if err!=nil{
			log.Warn(err)
		}else{
			excessThreshold=v
		}
	}

	return excessThreshold
}

func compactExcessThresholdValueSet(excessThreshold float64) error {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	err:=configureSvr.ConfigSet(CompactExcessThresholdValue,fmt.Sprintf("%v",excessThreshold))
	if err!=nil{
		log.Warn(err)
		return err
	}else{
		return nil
	}
}

func compactExcessThresholdMinCountGet() int64 {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	minCount:=int64(10000)
	value,err:=configureSvr.ConfigGet(CompactExcessThresholdMinCount)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			log.Debugf("compactExcessThresholdMinCount not exist,now set %v",minCount)
			err=configureSvr.ConfigSet(CompactExcessThresholdMinCount,fmt.Sprintf("%d",minCount))
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		v,err:=helper.StringToInt64(value)
		if err!=nil{
			log.Warn(err)
		}else{
			minCount=v
		}
	}

	return minCount
}

func compactExcessThresholdMinCountSet(minCount int64) error {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	err:=configureSvr.ConfigSet(CompactExcessThresholdMinCount,fmt.Sprintf("%d",minCount))
	if err!=nil{
		log.Warn(err)
		return err
	}else{
		return nil
	}
}
