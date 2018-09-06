package service

import (
	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeAddFree,ServiceVolumeAddFree)
}

func ServiceVolumeAddFree(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var count int = 1
	var bdir string
	var idir string

	value,ok:=m["count"]
	if ok {
		tmp,err:=helper.GetInt(value)
		if err==nil{
			count=int(tmp)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"count missing")
	}

	value,ok=m["bdir"]
	if ok {
		bdir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"bdir missing")
	}

	value,ok=m["idir"]
	if ok {
		idir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"idir missing")
	}

	sn,err:=storage.AddFreeVolume(count,bdir,idir)
	if err!=nil{
		log.Error(err)
		return helper.ResultBuildWithExtInfo(errors.RetVolumeAddFree,err.Error())
	}else{
		json:=gabs.New()
		json.Set(sn,"count")
		return helper.ResultBuildWithBody(constant.Success,json)
	}
}
