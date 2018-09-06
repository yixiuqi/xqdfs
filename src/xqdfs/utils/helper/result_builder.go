package helper

import (
	"github.com/Jeffail/gabs"
)

const (
	CmdCode  			= 	"code"
	CmdResult	 		= 	"result"
	CmdInfo			= 	"info"
	CmdSuccess			= 	"success"
)

func ResultBuild(result int32) *gabs.Container {
	json:=gabs.New()
	json.Set(float64(result),CmdResult)
	if result==0 {
		json.Set(CmdSuccess,CmdInfo)
	}else{
		json.Set(result,CmdInfo)
	}

	return json
}

func ResultBuildWithExtInfo(result int32,extInfo string) *gabs.Container {
	json:=gabs.New()
	json.Set(float64(result),CmdResult)
	if result==0 {
		if extInfo!=""{
			json.Set(extInfo,CmdInfo)
		}else {
			json.Set(CmdSuccess, CmdInfo)
		}
	}else{
		json.Set(extInfo,CmdInfo)
	}

	return json
}

func ResultBuildWithBody(result int32,json *gabs.Container) *gabs.Container {
	json.Set(float64(result),CmdResult)
	if json.ExistsP(CmdInfo)==false{
		if result == 0 {
			json.Set(CmdSuccess, CmdInfo)
		} else {
			json.Set(result, CmdInfo)
		}
	}

	return json
}



