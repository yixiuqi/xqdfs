package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeAdd,ServiceVolumeAdd)
}

func ServiceVolumeAdd(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var vid int32
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"vid missing")
	}

	_,err:=storage.AddVolume(vid)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetVolumeAdd,err.Error())
	}else{
		return helper.ResultBuild(constant.Success)
	}
}
