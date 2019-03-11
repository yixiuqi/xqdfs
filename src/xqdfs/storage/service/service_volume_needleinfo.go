package service

import (
	"bytes"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"
	
	"github.com/Jeffail/gabs"
	"github.com/json-iterator/go"
	jsonSys "encoding/json"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeNeedleInfo,ServiceVolumeNeedleInfo)
}

type RequestVolumeNeedleInfo struct {
	Vid int32 			`json:"vid"`
	Key int64 			`json:"key"`
}
func ServiceVolumeNeedleInfo(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeNeedleInfo{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	v:= storage.Volumes[req.Vid]
	if v != nil {
		n,err:= v.GetHeader(req.Key)
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
		}else{
			var b []byte
			b,err=json.Marshal(n)
			if err==nil{
				var item *gabs.Container
				dec := jsonSys.NewDecoder(bytes.NewBuffer(b))
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

