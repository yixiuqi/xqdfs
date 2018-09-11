package service

import (
	"bytes"

	"xqdfs/utils/helper"
	"xqdfs/storage/needle"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"
	"xqdfs/utils/plugin"
)

func init() {
	plugin.PluginAddService(constant.CmdVolumeUpload,ServiceVolumeUpload)
}

func ServiceVolumeUpload(m map[string]interface{}) interface{}{
	var storage *store.Store
	if s:=plugin.PluginGetObject(plugin.PlugineStorage);s==nil {
		log.Errorf("%s no support",plugin.PlugineStorage)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		storage=s.(*store.Store)
	}

	var replicationServer *replication.ReplicationServer
	if r:=plugin.PluginGetObject(plugin.PluginReplicationServer);r==nil {
		log.Errorf("%s no support",plugin.PluginReplicationServer)
		return helper.ResultBuild(errors.RetNoSupport)
	}else{
		replicationServer=r.(*replication.ReplicationServer)
	}

	var img []byte
	var vid int32
	var key int64
	var cookie int32
	value,ok:=m["img"]
	if ok {
		img=helper.ImageGet("",value.(string))
		if img==nil||len(img)==0 {
			return helper.ResultBuildWithExtInfo(errors.RetImageData,errors.ErrImageData.Error())
		}
	}else{
		return helper.ResultBuildWithExtInfo(errors.RetMissingParameter,"img missing")
	}

	value,ok=m["vid"]
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
		n:= needle.NewWriter(key, int32(cookie), int32(len(img)))
		defer n.Close()

		buf:= &bytes.Buffer{}
		buf.Write(img)

		err:= n.ReadFrom(buf)
		if err!=nil{
			log.Error(err)
			return helper.ResultBuildWithExtInfo(errors.RetOptUpload,err.Error())
		}

		err= v.Write(n)
		if err!=nil{
			if err==errors.ErrNeedleExist {
				replication,ok:=m["replication"]
				if ok && replication==true {
					p:=&process.ReplicationUpload{
						Vid:vid,
						Key:key,
						Cookie:cookie,
						Image:img,
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
			replication,ok:=m["replication"]
			if ok && replication==true {
				p:=&process.ReplicationUpload{
					Vid:vid,
					Key:key,
					Cookie:cookie,
					Image:img,
				}
				replicationServer.Replication(p)
			}
			return helper.ResultBuild(constant.Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
