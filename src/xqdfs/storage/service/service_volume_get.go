package service

import (
	"context"
	"net/http"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeGet,ServiceVolumeGet)
}

type RequestVolumeGet struct {
	Vid int32 			`json:"vid"`
	Key int64 			`json:"key"`
	Cookie int32 		`json:"cookie"`
	Replication bool 	`json:"replication"`
}
func ServiceVolumeGet(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeGet{}
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
		n,err:= v.Read(req.Key, req.Cookie)
		if err!=nil{
			log.Error(err)
			e,ok:=err.(errors.Error)
			if ok {
				return helper.ResultBuildWithExtInfo(int32(e),err.Error())
			}else{
				return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
			}
		}else{
			if inv.Method=="GET" {
				ctx.(*gin.Context).Data(http.StatusOK, "image/jpeg", n.Data)
				return nil
			}else{
				json:=gabs.New()
				json.Set(n.Data,"img")
				json.Set(req.Vid,"vid")
				json.Set(req.Key,"key")
				json.Set(req.Cookie,"cookie")
				return helper.ResultBuildWithBody(constant.Success,json)
			}
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
