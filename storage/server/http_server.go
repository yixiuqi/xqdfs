package server

import (
	"xqdfs/storage/conf"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"fmt"
	"xqdfs/utils/log"
	"xqdfs/storage/store"
	"os"
	"xqdfs/storage/service"
)

type HttpServer struct {
}

func NewHttpServer(c *conf.Config,s *store.Store) (server *HttpServer,err error) {
	server=new(HttpServer)
	err=server.Start(c,s)
	return
}

func (this *HttpServer) Start(c *conf.Config,s *store.Store) (err error){
	go this.process(c,s)
	return
}

func (this *HttpServer) process(c *conf.Config,s *store.Store) {
	context:=service.NewContext(c,s)
	router:= gin.Default()
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	router.POST(service.HttpVolumeAddFree, NewHttpWrap(context,service.ServiceVolumeAddFree).Handler)
	router.POST(service.HttpVolumeAdd, NewHttpWrap(context,service.ServiceVolumeAdd).Handler)
	router.POST(service.HttpVolumeClear, NewHttpWrap(context,service.ServiceVolumeClear).Handler)
	router.POST(service.HttpVolumeCompact, NewHttpWrap(context,service.ServiceVolumeCompact).Handler)
	router.POST(service.HttpVolumeUpload, NewHttpWrap(context,service.ServiceVolumeUpload).Handler)
	router.POST(service.HttpVolumeGet, NewHttpWrap(context,service.ServiceVolumeGet).Handler)
	router.GET(service.HttpVolumeGet, NewHttpWrap(context,service.ServiceVolumeGet).Handler)
	router.POST(service.HttpVolumeDelete, NewHttpWrap(context,service.ServiceVolumeDelete).Handler)
	router.POST(service.HttpVolumeNeedleInfo, NewHttpWrap(context,service.ServiceVolumeNeedleInfo).Handler)

	router.POST(service.HttpStoreStat, NewHttpWrap(context,service.ServiceStoreStat).Handler)
	router.POST(service.HttpStoreConf, NewHttpWrap(context,service.ServiceStoreConf).Handler)
	router.POST(service.HttpStoreUpload, NewHttpWrap(context,service.ServiceStoreUpload).Handler)
	router.POST(service.HttpStoreGet, NewHttpWrap(context,service.ServiceStoreGet).Handler)
	router.GET(service.HttpStoreGet, NewHttpWrap(context,service.ServiceStoreGet).Handler)

	err:= router.Run(fmt.Sprintf(":%d",context.Conf.Http.Port))
	if err != nil {
		log.Errorf("http error[%v]",err)
		os.Exit(1)
	}
}

func (this *HttpServer) Stop() {
}


