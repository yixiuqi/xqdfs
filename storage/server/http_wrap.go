package server

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"xqdfs/utils/helper"
	"net/http"
	"github.com/Jeffail/gabs"
	"xqdfs/utils/log"
	"strings"
	"xqdfs/storage/service"
)

type HttpWrap struct {
	context *service.Context
	handle service.HandlerFunc
}

func NewHttpWrap(context *service.Context,handle service.HandlerFunc) *HttpWrap {
	item:=new(HttpWrap)
	item.context=context
	item.handle=handle
	return item
}

func (this *HttpWrap) Handler(c *gin.Context) {
	helper.HandleErr()
	contentType := c.GetHeader("Content-Type")
	method:=c.Request.Method

	m:= make(map[string]interface{})
	if strings.Contains(contentType, "text/plain") || strings.Contains(contentType,"application/json"){
		body := new(bytes.Buffer)
		buf := make([]byte, 1024*128)
		n, err:= c.Request.Body.Read(buf)
		for n > 0 {
			body.Write(buf[0:n])
			n, err = c.Request.Body.Read(buf)
		}

		dec := json.NewDecoder(bytes.NewBuffer(body.Bytes()))
		dec.UseNumber()
		err = dec.Decode(&m)
		if err != nil {
			debug:=helper.NewStringBuilder().Append("http body error[").Append(err).Append("]").ToString()
			log.Error(debug)
			result := helper.ResultBuildWithExtInfo(service.Failed,debug)
			c.JSON(http.StatusOK, result)
		}
	}else{
		query := c.Request.URL.Query()
		for k, v := range query {
			m[k] = v[0]
		}

		req := c.Request
		req.ParseForm()
		req.ParseMultipartForm(32 << 20)
		form := req.PostForm
		for k, v := range form {
			m[k] = v[0]
		}
	}

	m["http_contentType"]=contentType
	m["http_method"]=method
	m["http_context"]=c
	ret:=this.handle(this.context,m)
	if ret!=nil{
		json:=ret.(*gabs.Container)
		c.JSON(http.StatusOK, json.Data())
	}
}
