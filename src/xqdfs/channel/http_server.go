package channel

import (
	"fmt"
	"os"
	"net/http"
	"encoding/base64"

	"xqdfs/utils/log"
	"xqdfs/utils/plugin"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/Jeffail/gabs"
)

type HttpServer struct {
}

func NewHttpServer(port int) (*HttpServer,error) {
	server:=&HttpServer{}
	go server.process(port)
	return server,nil
}

func (this *HttpServer) process(port int) {
	router:= gin.Default()
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)
	router.Static("/page", "webroot")

	services:=plugin.PluginGetServices()
	for k,v:=range services {
		log.Debugf("path[%s]",k)
		router.GET(k, NewHttpWrap(v).Handler)
		router.POST(k, NewHttpWrap(v).Handler)
	}

	router.POST("/system/upload", this.uploadFile)
	err:= router.Run(fmt.Sprintf(":%d",port))
	if err != nil {
		log.Errorf("http error[%v]",err)
		os.Exit(1)
	}
}

func (this HttpServer) uploadFile(c *gin.Context) {
	var defaultMaxMemory int64 = 32 << 20 // 32 MB
	c.Request.ParseMultipartForm(defaultMaxMemory)
	file := c.Request.MultipartForm.File
	for _, v := range file {
		reader, _ := v[0].Open()
		defer reader.Close()
		var buf []byte = make([]byte, defaultMaxMemory)
		len, _ := reader.Read(buf)
		ct := v[0].Header["Content-Type"][0]

		json := gabs.New()
		str := "data:" + ct + ";base64," + base64.StdEncoding.EncodeToString(buf[0:len])
		json.Set(0, "result")
		json.Set(str, "image")
		c.JSON(http.StatusOK, json.Data())
		return
	}
}

func (this *HttpServer) Stop() {
	log.Info("HttpServer stop")
}

