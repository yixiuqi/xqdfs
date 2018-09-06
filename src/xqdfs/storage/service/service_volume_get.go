package service

import (
	"net/http"

	"xqdfs/errors"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/utils/helper"

	"github.com/gin-gonic/gin"
	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeGet,ServiceVolumeGet)
}

func ServiceVolumeGet(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var vid int32
	var key int64
	var cookie int32
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

	value,ok=m["cookie"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			cookie=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"cookie missing")
	}

	v:= storage.Volumes[vid]
	if v != nil {
		n,err:= v.Read(key, cookie)
		if err!=nil{
			log.Error(err)
			e,ok:=err.(errors.Error)
			if ok {
				return helper.ResultBuildWithExtInfo(int32(e),err.Error())
			}else{
				return helper.ResultBuildWithExtInfo(errors.RetOptGet,err.Error())
			}
		}else{
			c,ok:=m["http_context"]
			if ok {
				method,ok:=m["http_method"]
				if ok {
					if method=="GET" {
						c.(*gin.Context).Data(http.StatusOK, "image/jpeg", n.Data)
						return nil
					}
				}
			}

			json:=gabs.New()
			json.Set(n.Data,"img")
			json.Set(vid,"vid")
			json.Set(key,"key")
			json.Set(cookie,"cookie")
			return helper.ResultBuildWithBody(constant.Success,json)
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
