package alloc

import (
	"fmt"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/master/configure"
)

const(
	AllocOrderMinFreeSpace 	= 	"AllocOrderMinFreeSpace"			// 每个卷预留多少空间 default:104857600
	AllocOrderConsumeCount 	= 	"AllocOrderConsumeCount"			// 同时default:3
)

var(
	configureSvr *configure.ConfigureServer
)

func setupConfigureSvr(conf *configure.ConfigureServer) {
	configureSvr=conf
}

func allocOrderMinFreeSpaceGet() int64 {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	var orderMinFreeSpace int64=104857600
	value,err:=configureSvr.ConfigGet(AllocOrderMinFreeSpace)
	if err!=nil {
		if err==errors.ErrParamNotExist {
			log.Debugf("orderMinFreeSpace not exist,now set %d",orderMinFreeSpace)
			err=configureSvr.ConfigSet(AllocOrderMinFreeSpace,fmt.Sprintf("%d",orderMinFreeSpace))
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		v,err:=helper.StringToInt64(value)
		if err!=nil{
			log.Warn(err)
		}else{
			orderMinFreeSpace=v
		}
	}
	return orderMinFreeSpace
}

func allocOrderMinFreeSpaceSet(orderMinFreeSpace int64) error {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	err:=configureSvr.ConfigSet(AllocOrderMinFreeSpace,fmt.Sprintf("%d",orderMinFreeSpace))
	if err!=nil{
		log.Warn(err)
		return err
	}else{
		return nil
	}
}

func allocOrderConsumeCountGet() int {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	orderConsumeCount:=3
	value,err:=configureSvr.ConfigGet(AllocOrderConsumeCount)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			log.Debugf("orderConsumeCount not exist,now set %d",orderConsumeCount)
			err=configureSvr.ConfigSet(AllocOrderConsumeCount,fmt.Sprintf("%d",orderConsumeCount))
			if err!=nil{
				log.Warn(err)
			}
		}
	}else{
		v,err:=helper.StringToInt(value)
		if err!=nil{
			log.Warn(err)
		}else{
			orderConsumeCount=v
		}
	}
	return orderConsumeCount
}

func allocOrderConsumeCountSet(orderConsumeCount int) error {
	if configureSvr==nil {
		panic("configureSvr is null")
	}

	err:=configureSvr.ConfigSet(AllocOrderConsumeCount,fmt.Sprintf("%d",orderConsumeCount))
	if err!=nil{
		log.Warn(err)
		return err
	}else{
		return nil
	}
}

