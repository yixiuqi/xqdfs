package service

import (
	"fmt"
	"os"
	"encoding/base64"
	"net/http"

	"xqdfs/utils/log"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/proxy"
	"xqdfs/constant"
	"xqdfs/master/conf"
	"xqdfs/master/strategy"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/Jeffail/gabs"
)

type HttpServer struct {
}

func NewHttpServer(conf *conf.Config,configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,strategyServer *strategy.AllocStrategyServer,proxyStorage *proxy.ProxyStorage) (server *HttpServer,err error) {
	context:=NewContext(conf,configureServer,discoveryServer,strategyServer,proxyStorage)
	server=new(HttpServer)
	err=server.Start(context)
	return
}

func (this *HttpServer) Start(context *Context) (err error){
	go this.process(context)
	return
}

func (this *HttpServer) process(context *Context) {
	router:= gin.Default()
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)
	router.Static("/page", "webroot")

	router.POST(constant.HttpGroupGetAll, NewHttpWrap(context,ServiceGroupGetAll).Handler)
	router.POST(constant.HttpGroupReadOnly, NewHttpWrap(context,ServiceGroupReadOnly).Handler)
	router.POST(constant.HttpGroupAdd, NewHttpWrap(context,ServiceGroupAdd).Handler)
	router.POST(constant.HttpGroupRemove, NewHttpWrap(context,ServiceGroupRemove).Handler)
	router.POST(constant.HttpGroupAddStorage, NewHttpWrap(context,ServiceGroupAddStorage).Handler)
	router.POST(constant.HttpGroupRemoveStorage, NewHttpWrap(context,ServiceGroupRemoveStorage).Handler)

	router.POST(constant.HttpStorageVolumeGetAll, NewHttpWrap(context,ServiceStorageVolumeGetAll).Handler)
	router.POST(constant.HttpStorageVolumeClear, NewHttpWrap(context,ServiceStorageVolumeClear).Handler)
	router.POST(constant.HttpStorageVolumeCompact, NewHttpWrap(context,ServiceStorageVolumeCompact).Handler)
	router.POST(constant.HttpStorageInit, NewHttpWrap(context,ServiceStorageInit).Handler)
	router.POST(constant.HttpStorageAdd, NewHttpWrap(context,ServiceStorageAdd).Handler)
	router.POST(constant.HttpStorageRemove, NewHttpWrap(context,ServiceStorageRemove).Handler)
	router.POST(constant.HttpStorageGetAll, NewHttpWrap(context,ServiceStorageGetAll).Handler)
	router.POST(constant.HttpStorageGetConfigure, NewHttpWrap(context,ServiceStorageGetConfigure).Handler)

	router.POST(constant.HttpUsageGroups, NewHttpWrap(context,ServiceUsageGroups).Handler)

	router.POST(constant.HttpOptUpload, NewHttpWrap(context,ServiceOptUpload).Handler)
	router.POST(constant.HttpOptGet, NewHttpWrap(context,ServiceOptGet).Handler)
	router.POST(constant.HttpOptDelete, NewHttpWrap(context,ServiceOptDelete).Handler)

	router.POST("/system/upload", this.uploadFile)
	err:= router.Run(fmt.Sprintf(":%d",context.Conf.Http.Port))
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

