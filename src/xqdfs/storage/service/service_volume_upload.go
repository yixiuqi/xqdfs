package service

import (
	"bytes"
	"context"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/storage/store"
	"xqdfs/storage/needle"
	"xqdfs/storage/replication"
	"xqdfs/storage/replication/process"

	"github.com/json-iterator/go"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeUpload,ServiceVolumeUpload)
}

type RequestVolumeUpload struct {
	Vid int32 			`json:"vid"`
	Key int64 			`json:"key"`
	Cookie int32 		`json:"cookie"`
	Img []byte 		`json:"img"`
	Replication bool 	`json:"replication"`
}
func ServiceVolumeUpload(ctx context.Context,inv *plugin.Invocation) interface{}{
	req:=&RequestVolumeUpload{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err:=json.Unmarshal(inv.Body,req)
	if err!=nil {
		log.Warn(err)
		return helper.ResultBuildWithExtInfo(errors.RetParameterError,err.Error())
	}

	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Warnf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var replicationServer *replication.ReplicationServer
	if r:=plugin.PluginGetObject(plugin.PluginReplicationServer);r==nil {
		log.Warnf("%s no support",plugin.PluginReplicationServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		replicationServer=r.(*replication.ReplicationServer)
	}

	v:= storage.Volumes[req.Vid]
	if v != nil {
		n:= needle.NewWriter(req.Key, req.Cookie, int32(len(req.Img)))
		defer n.Close()

		buf:= &bytes.Buffer{}
		buf.Write(req.Img)

		err:= n.ReadFrom(buf)
		if err!=nil{
			log.Warn(err)
			return helper.ResultBuildWithExtInfo(errors.RetOptUpload,err.Error())
		}

		err= v.Write(n)
		if err!=nil{
			if err==errors.ErrNeedleExist {
				if req.Replication==true {
					p:=&process.ReplicationUpload{
						Vid:req.Vid,
						Key:req.Key,
						Cookie:req.Cookie,
						Image:req.Img,
					}
					replicationServer.Replication(p)
				}
			}

			e,ok:=err.(errors.Error)
			if ok {
				return helper.ResultBuildWithExtInfo(int32(e),err.Error())
			}else{
				return helper.ResultBuildWithExtInfo(errors.RetOptUpload,err.Error())
			}
		}else{
			if req.Replication==true {
				p:=&process.ReplicationUpload{
					Vid:req.Vid,
					Key:req.Key,
					Cookie:req.Cookie,
					Image:req.Img,
				}
				replicationServer.Replication(p)
			}
			return helper.ResultBuild(constant.Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
