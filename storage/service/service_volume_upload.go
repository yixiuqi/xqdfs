package service

import (
	"bytes"

	"xqdfs/utils/helper"
	"xqdfs/storage/needle"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/replication/process"
)

/**
 * @api {post} /volume/upload [Volume]图片上传
 * @apiDescription [Volume]图片上传
 * @apiGroup Storage
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {string} img 图片数据
 * @apiParam {int32} vid volume id
 * @apiParam {int64} key 文件key
 * @apiParam {int32} cookie 文件cookie
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/upload
{
	"seq": "f428b43e-5bc1-4bea-a93a-ed177ca4552e"
}
 * @apiErrorExample 失败返回示例
{
    "info": "查询错误",
    "result": 1
}

 * @apiSuccessExample 成功返回示例
{
    "info": "success",
    "result": 0
}
* */
func ServiceVolumeUpload(context *Context,m map[string]interface{}) interface{}{
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

	log.Debugf("upload volume[%d] key[%d] cookie[%d]",vid,key,cookie)

	v:= context.Store.Volumes[vid]
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
			log.Error(err)
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
				context.ReplicationServer.Replication(p)
			}else{
				log.Debugf("receive replication request vid[%d] key[%d]",vid,key)
			}
			return helper.ResultBuild(constant.Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(errors.RetVolumeNotExist,errors.ErrVolumeNotExist.Error())
	}
}
