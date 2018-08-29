package service

import (
	"xqdfs/utils/helper"
	"github.com/Jeffail/gabs"
	"bytes"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"xqdfs/errors"
	"xqdfs/storage/needle"
	"net/http"
	"xqdfs/utils/log"
)

/**
 * @api {post} /volume/addfree 创建空闲块
 * @apiDescription 创建空闲块
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
* @apiParam {int} count 创建数目
* @apiParam {string} bdir block文件目录
* @apiParam {string} idir block.idx文件目录
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/addfree
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
func ServiceVolumeAddFree(context *Context,m map[string]interface{}) interface{}{
	var count int = 1
	var bdir string
	var idir string

	value,ok:=m["count"]
	if ok {
		tmp,err:=helper.GetInt(value)
		if err==nil{
			count=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing count")
	}

	value,ok=m["bdir"]
	if ok {
		bdir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing bdir")
	}

	value,ok=m["idir"]
	if ok {
		idir=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing idir")
	}

	sn,err:=context.Store.AddFreeVolume(count,bdir,idir)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(Failed,err.Error())
	}else{
		json:=gabs.New()
		json.Set(sn,"count")
		return helper.ResultBuildWithBody(Success,json)
	}
}

/**
 * @api {post} /volume/add 块启用
 * @apiDescription 块启用
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
* @apiParam {int} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/add
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
func ServiceVolumeAdd(context *Context,m map[string]interface{}) interface{}{
	var vid int32
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	_,err:=context.Store.AddVolume(vid)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(Failed,err.Error())
	}else{
		return helper.ResultBuild(Success)
	}
}

/**
 * @api {post} /volume/clear 块清空
 * @apiDescription 块清空
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
* @apiParam {int} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/clear
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
func ServiceVolumeClear(context *Context,m map[string]interface{}) interface{}{
	var vid int32
	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	v:= context.Store.Volumes[vid]
	if v != nil {
		err:= v.Clear()
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}else{
			return helper.ResultBuild(Success)
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,errors.ErrVolumeNotExist.Error())
	}
}

/**
 * @api {post} /volume/upload 图片上传
 * @apiDescription 图片上传
 * @apiGroup Volume
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
			return helper.ResultBuildWithExtInfo(Failed,"image data error")
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing img")
	}

	value,ok=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing key")
	}

	value,ok=m["cookie"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			cookie=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing cookie")
	}

	v:= context.Store.Volumes[vid]
	if v != nil {
		n:= needle.NewWriter(key, int32(cookie), int32(len(img)))
		defer n.Close()

		buf:= &bytes.Buffer{}
		buf.Write(img)

		err:= n.ReadFrom(buf)
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}

		err= v.Write(n)
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}else{
			return helper.ResultBuild(Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(Failed,errors.ErrVolumeNotExist.Error())
	}
}

/**
 * @api {post} /volume/get 图片下载
 * @apiDescription 图片下载
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiParam {int64} key 文件key
 * @apiParam {int32} cookie 文件cookie
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (成功返回参数) {string} img 图片
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/get
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
func ServiceVolumeGet(context *Context,m map[string]interface{}) interface{}{
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
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing key")
	}

	value,ok=m["cookie"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			cookie=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing cookie")
	}

	v:= context.Store.Volumes[vid]
	if v != nil {
		n,err:= v.Read(key, int32(cookie))
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
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
			return helper.ResultBuildWithBody(Success,json)
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,errors.ErrVolumeNotExist.Error())
	}
}

/**
 * @api {post} /volume/delete 图片删除
 * @apiDescription 图片删除
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiParam {int64} key 文件key
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/delete
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
func ServiceVolumeDelete(context *Context,m map[string]interface{}) interface{}{
	var vid int32
	var key int64

	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing key")
	}

	v:= context.Store.Volumes[vid]
	if v != nil {
		err:= v.Delete(key)
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}else{
			return helper.ResultBuild(Success)
		}
	} else {
		return helper.ResultBuildWithExtInfo(Failed,errors.ErrVolumeNotExist.Error())
	}
}

/**
 * @api {post} /volume/needleinfo 查询Needle信息
 * @apiDescription 查询Needle信息
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiParam {int64} key 文件key
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/needleinfo
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
func ServiceVolumeNeedleInfo(context *Context,m map[string]interface{}) interface{}{
	var vid int32
	var key int64

	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	value,ok=m["key"]
	if ok {
		tmp,err:=helper.GetInt64(value)
		if err==nil{
			key=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing key")
	}

	v:= context.Store.Volumes[vid]
	if v != nil {
		n,err:= v.GetHeader(key)
		if err!=nil{
			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}else{
			var b []byte
			b,err=json.Marshal(n)
			if err==nil{
				var item *gabs.Container
				item,err=gabs.ParseJSON(b)
				if err==nil{
					json:=gabs.New()
					json.Set(item.Data(),"needle")
					return helper.ResultBuildWithBody(Success,json)
				}
			}

			return helper.ResultBuildWithExtInfo(Failed,err.Error())
		}
	} else {
		return helper.ResultBuildWithExtInfo(Failed,errors.ErrVolumeNotExist.Error())
	}
}

/**
 * @api {post} /volume/compact 块回收
 * @apiDescription 块回收
 * @apiGroup Volume
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {int32} vid volume id
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/volume/compact
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
func ServiceVolumeCompact(context *Context,m map[string]interface{}) interface{}{
	var vid int32

	value,ok:=m["vid"]
	if ok {
		tmp,err:=helper.GetInt32(value)
		if err==nil{
			vid=tmp
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing vid")
	}

	go func() {
		err:=context.Store.CompactVolume(vid)
		if err!=nil{
			log.Infof("volume[%v] compact error",err)
		}else{
			log.Infof("volume[%v] compact ok",vid)
		}
	}()
	return helper.ResultBuild(Success)
}