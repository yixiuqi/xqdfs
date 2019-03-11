package channel

import (
	"sync"
	"bytes"
	"strings"
	"net/http"
	"encoding/json"

	"xqdfs/errors"
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

var BufPool = sync.Pool{
	New: func() interface{} { return make([]byte, 1024*32) },
}
var InvPool = sync.Pool{
	New: func() interface{} { return &plugin.Invocation{} },
}
func (this *HttpWrap) Handler(c *gin.Context) {
	this.wg.Add(1)
	defer func(){
		this.wg.Done()
	}()
	defer helper.HandleErr()
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "text/plain") || strings.Contains(contentType,"application/json"){
		if strings.Contains(contentType, plugin.HttpTextPlain) {
			contentType=plugin.HttpTextPlain
		}
		if strings.Contains(contentType,plugin.HttpApplicationJson) {
			contentType=plugin.HttpApplicationJson
		}

		body := new(bytes.Buffer)
		//pool
		buf:= BufPool.Get().([]byte)
		defer BufPool.Put(buf)
		n,_:= c.Request.Body.Read(buf)
		for n > 0 {
			body.Write(buf[0:n])
			n,_= c.Request.Body.Read(buf)
		}

		inv:=InvPool.Get().(*plugin.Invocation)
		defer InvPool.Put(inv)

		inv.Body=body.Bytes()
		inv.Method=c.Request.Method
		inv.ContentType=contentType
		result:= this.handle(c,inv)
		if result!=nil{
			ret:=result.(*gabs.Container)
			c.JSON(http.StatusOK, ret.Data())
		}
	}else{
		m:= make(map[string]interface{})
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

		b, err:= json.Marshal(m)
		if err!=nil{
			c.JSON(http.StatusOK, helper.ResultBuild(errors.RetParameterError).Data())
			return
		}
		inv:=&plugin.Invocation{}
		inv.Body=b
		inv.Method=c.Request.Method
		inv.ContentType=contentType
		result:= this.handle(c,inv)
		if result!=nil{
			ret:=result.(*gabs.Container)
			c.JSON(http.StatusOK, ret.Data())
		}
	}
}
