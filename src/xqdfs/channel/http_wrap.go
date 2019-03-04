package channel

import (
	"sync"
	"bytes"
	"strings"
	"net/http"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"

	"github.com/Jeffail/gabs"
	"github.com/gin-gonic/gin"
)

type HttpWrap struct {
	wg *sync.WaitGroup
	handle plugin.HandlerFunc
}

func NewHttpWrap(wg *sync.WaitGroup,handle plugin.HandlerFunc) *HttpWrap {
	wrap:=&HttpWrap{
		wg:wg,
		handle:handle,
	}
	return wrap
}

func (this *HttpWrap) Handler(c *gin.Context) {
	this.wg.Add(1)
	defer func(){
		this.wg.Done()
	}()
	defer helper.HandleErr()
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
			result := helper.ResultBuildWithExtInfo(errors.RetParameterError,debug)
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
	ret:=this.handle(m)
	if ret!=nil{
		json:=ret.(*gabs.Container)
		c.JSON(http.StatusOK, json.Data())
	}
}
