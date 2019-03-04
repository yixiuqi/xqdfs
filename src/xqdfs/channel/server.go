package channel

import (
	"xqdfs/utils/log"
	"xqdfs/utils/conf"
)

type Server struct {
	httpServer *HttpServer
}

func NewServer(conf *conf.Server) (server *Server,err error) {
	var(
		httpServer *HttpServer
	)

	if httpServer, err = NewHttpServer(conf.Port); err != nil {
		log.Errorf("http server init error[%v]",err)
		return
	}

	server=&Server{
		httpServer:httpServer,
	}
	return
}

func (this *Server) Stop() {
	log.Info("Server stop")

	if this.httpServer !=nil {
		this.httpServer.Stop()
	}
}
