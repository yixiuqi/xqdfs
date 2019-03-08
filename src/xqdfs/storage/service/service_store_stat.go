package service

import (
	"sort"
	"bytes"
	"context"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"
	"xqdfs/utils/plugin"
	"xqdfs/storage/volume"

	"github.com/Jeffail/gabs"
)

func init() {
	plugin.PluginAddService(constant.CmdStoreStat,ServiceStoreStat)
}

type VolumeSortById []*volume.Volume
func (v VolumeSortById) Len() int {
	return len(v)
}
func (v VolumeSortById) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v VolumeSortById) Less(i, j int) bool {
	return v[j].Id > v[i].Id
}

/**
 * @api {post} /store/stat [Store]查询状态
 * @apiDescription [Store]查询当前存储节点状态信息
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 非0错误码
 * @apiError (失败返回参数) {string} info 信息
* */
func ServiceStoreStat(ctx context.Context,inv *plugin.Invocation) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	jsonStat:=gabs.New()
	jsonStat.Set(storage.FreeId,"freeId")

	jsonStat.Array("volumes")
	vs:=make([]*volume.Volume,0)
	for _,v:=range storage.Volumes{
		vs=append(vs,v)
	}
	sort.Sort(VolumeSortById(vs))

	for _,v:=range vs{
		b,err:=json.Marshal(v)
		if err==nil{
			dec := json.NewDecoder(bytes.NewBuffer(b))
			dec.UseNumber()
			item,err:=gabs.ParseJSONDecoder(dec)
			item.Set(v.ImageCount(),"imageCount")
			if err==nil {
				jsonStat.ArrayAppend(item.Data(), "volumes")
			}
		}else{
			log.Error(err)
		}
	}

	jsonStat.Array("freeVolumes")
	for _,v:=range storage.FreeVolumes{
		b,err:=json.Marshal(v)
		if err==nil{
			dec := json.NewDecoder(bytes.NewBuffer(b))
			dec.UseNumber()
			item,err:=gabs.ParseJSONDecoder(dec)
			if err==nil{
				jsonStat.ArrayAppend(item.Data(),"freeVolumes")
			}
		}else{
			log.Error(err)
		}
	}

	stat,err:=json.Marshal(storage.Stats)
	if err==nil{
		dec := json.NewDecoder(bytes.NewBuffer(stat))
		dec.UseNumber()
		item,err:=gabs.ParseJSONDecoder(dec)
		if err==nil{
			jsonStat.Set(item.Data(),"stats")
		}
	}else{
		log.Error(err)
	}

	result:=helper.ResultBuildWithBody(constant.Success,jsonStat)
	return result
}
