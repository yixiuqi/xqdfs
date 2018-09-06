package service

import (
	"encoding/json"
	"bytes"

	"xqdfs/utils/helper"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
	
	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeNeedleInfo,ServiceVolumeNeedleInfo)
}

func ServiceVolumeNeedleInfo(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var vid int32
	var key int64
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"vid missing")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"key missing")
	}

	v:= storage.Volumes[vid]
	if v != nil {
		n,err:= v.GetHeader(key)
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
		}else{
			var b []byte
			b,err=json.Marshal(n)
			if err==nil{
				var item *gabs.Container
				dec := json.NewDecoder(bytes.NewBuffer(b))
				dec.UseNumber()
				item,err=gabs.ParseJSONDecoder(dec)
				if err==nil{
					json:=gabs.New()
					json.Set(item.Data(),"needle")
					return helper.ResultBuildWithBody(constant.Success,json)
				}
			}

			return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}

