package service

import (
	"fmt"
	"os"

	"xqdfs/utils/log"
	"xqdfs/constant"
	"xqdfs/storage/conf"
	"xqdfs/storage/store"
	"xqdfs/storage/replication"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

type HttpServer struct {
}

func NewHttpServer(c *conf.Config,s *store.Store,replicationServer *replication.ReplicationServer) (server *HttpServer,err error) {
	server=new(HttpServer)
	err=server.Start(c,s,replicationServer)
	return
}

func (this *HttpServer) Start(c *conf.Config,s *store.Store,replicationServer *replication.ReplicationServer) (err error){
	go this.process(c,s,replicationServer)
	return
}

func (this *HttpServer) process(c *conf.Config,s *store.Store,replicationServer *replication.ReplicationServer) {
	context:=NewContext(c,s,replicationServer)
	router:= gin.Default()
	router.Use(cors.Default())
	gin.SetMode(gin.ReleaseMode)

	router.POST(constant.HttpVolumeAddFree, NewHttpWrap(context,ServiceVolumeAddFree).Handler)
	router.POST(constant.HttpVolumeAdd, NewHttpWrap(context,ServiceVolumeAdd).Handler)
	router.POST(constant.HttpVolumeClear, NewHttpWrap(context,ServiceVolumeClear).Handler)
	router.POST(constant.HttpVolumeCompact, NewHttpWrap(context,ServiceVolumeCompact).Handler)
	router.POST(constant.HttpVolumeCompactStatus, NewHttpWrap(context,ServiceVolumeCompactStatus).Handler)
	router.POST(constant.HttpVolumeUpload, NewHttpWrap(context,ServiceVolumeUpload).Handler)
	router.POST(constant.HttpVolumeGet, NewHttpWrap(context,ServiceVolumeGet).Handler)
	router.GET(constant.HttpVolumeGet, NewHttpWrap(context,ServiceVolumeGet).Handler)
	router.POST(constant.HttpVolumeDelete, NewHttpWrap(context,ServiceVolumeDelete).Handler)
	router.POST(constant.HttpVolumeNeedleInfo, NewHttpWrap(context,ServiceVolumeNeedleInfo).Handler)

	router.POST(constant.HttpStoreInit, NewHttpWrap(context,ServiceStoreInit).Handler)
	router.POST(constant.HttpStoreStat, NewHttpWrap(context,ServiceStoreStat).Handler)
	router.POST(constant.HttpStoreConf, NewHttpWrap(context,ServiceStoreConf).Handler)

	err:= router.Run(fmt.Sprintf(":%d",context.Conf.Http.Port))
	if err != nil {
		log.Errorf("http error[%v]",err)
		os.Exit(1)
	}
}

func (this *HttpServer) Stop() {
	log.Info("HttpServer stop")
}


