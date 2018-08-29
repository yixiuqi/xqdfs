package service

import (
	"bytes"
	"xqdfs/utils/helper"
	"xqdfs/storage/needle"
	"github.com/Jeffail/gabs"
	"fmt"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
	"xqdfs/errors"
)

/**
 * @api {post} /store/upload 图片上传
 * @apiDescription 图片上传
 * @apiGroup Store
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {string} img 图片数据
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/store/upload
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
func ServiceStoreUpload(context *Context,m map[string]interface{}) interface{}{
	var img []byte
	value,ok:=m["img"]
	if ok {
		img=helper.ImageGet("",value.(string))
		if img==nil||len(img)==0 {
			return helper.ResultBuildWithExtInfo(Failed,"image data error")
		}
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing img")
	}

	key:=helper.KeyGenerate()
	cookie:=Cookie

	n:= needle.NewWriter(key, cookie, int32(len(img)))
	defer n.Close()

	buf:= &bytes.Buffer{}
	buf.Write(img)

	err:= n.ReadFrom(buf)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(Failed,err.Error())
	}

	vid,err:=context.Store.Write(n)
	if err==nil{
		json:=gabs.New()
		json.Set(helper.KeyGetTime(key),"time")
		url:=fmt.Sprintf("%v,%v",vid,key)
		json.Set(url,"url_src")
		encodeString := base64.StdEncoding.EncodeToString([]byte(url))
		json.Set(encodeString,"url")
		return helper.ResultBuildWithBody(Success,json)
	} else {
		return helper.ResultBuildWithExtInfo(Failed,err.Error())
	}
}

/**
 * @api {post} /store/get 图片下载
 * @apiDescription 图片下载
 * @apiGroup Store
 * @apiVersion 1.0.0
 * @apiParam {string} [seq] 会话序号(非必填)
 * @apiParam {string} url 图片路径
 * @apiSuccess (成功返回参数) {int32} result 0表示成功
 * @apiError (失败返回参数) {int32} result 1表示失败
 * @apiError (失败返回参数) {string} info 信息
 * @apiExample 使用示例:
 * http://ip:port/store/get
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
func ServiceStoreGet(context *Context,m map[string]interface{}) interface{}{
	var url string
	value,ok:=m["url"]
	if ok {
		url=value.(string)
	}else{
		return helper.ResultBuildWithExtInfo(Failed,"missing url")
	}

	cookie:=Cookie
	dec,err:=base64.StdEncoding.DecodeString(url)
	if err!=nil{
		return helper.ResultBuildWithExtInfo(Failed,err.Error())
	}
	item:=strings.Split(string(dec),",")
	if len(item)!=2 {
		return helper.ResultBuildWithExtInfo(Failed,"url error:"+string(dec))
	}

	vid,err:=helper.StringToInt32(item[0])
	if err!=nil {
		return helper.ResultBuildWithExtInfo(Failed,"url[vid] error:"+string(dec))
	}

	key,err:=helper.StringToInt64(item[1])
	if err!=nil {
		return helper.ResultBuildWithExtInfo(Failed,"url[key] error")
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

