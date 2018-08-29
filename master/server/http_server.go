package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"fmt"
	"xqdfs/utils/log"
	"os"
	"bytes"
	"encoding/json"
	"xqdfs/utils/helper"
	"net/http"
	"errors"
	"xqdfs/master/conf"
	"xqdfs/master/service"
)

type HttpServer struct {
}

func NewHttpServer(c *conf.Config) (server *HttpServer,err error) {
	server=new(HttpServer)
	err=server.Start(c)
	return
}

func (this *HttpServer) Start(c *conf.Config) (err error){
	go this.process(c)
	return
}

func (this *HttpServer) process(c *conf.Config) {
	context:=service.NewContext(c)
	router:= gin.Default()
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	err:= router.Run(fmt.Sprintf(":%d",context.Conf.Http.Port))
	if err != nil {
		log.Errorf("http error[%v]",err)
		os.Exit(1)
	}
}

func (this *HttpServer) Stop() {
}

func (this *HttpServer) handle(c *gin.Context) (m map[string]interface{},err error){
	body := new(bytes.Buffer)
	buf := make([]byte, 1024*128)
	var n int
	n, err = c.Request.Body.Read(buf)
	for n > 0 {
		body.Write(buf[0:n])
		n, err = c.Request.Body.Read(buf)
	}

	m = make(map[string]interface{})
	dec := json.NewDecoder(bytes.NewBuffer(body.Bytes()))
	dec.UseNumber()
	err = dec.Decode(&m)

	if err != nil {
		debug:=helper.NewStringBuilder().Append("http body error[").Append(err).Append("]").ToString()
		err=errors.New(debug)

		result := helper.ResultBuildWithExtInfo(service.Failed,err.Error())
		c.JSON(http.StatusOK, result)
	}
	return
}
