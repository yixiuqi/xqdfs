package channel

import (
	"xqdfs/utils/log"
	"xqdfs/utils/conf"
)

type Server struct {
	httpServer *HttpServer
	thriftServer *ThriftServer
}

func NewServer(conf *conf.Server) (server *Server,err error) {
	var(
		httpServer *HttpServer
		thriftServer *ThriftServer
	)

	if httpServer, err = NewHttpServer(conf.Port); err != nil {
		log.Errorf("http server init error[%v]",err)
		return
	}

	//if thriftServer, err = NewThriftServer(conf.Port); err != nil {
	//	log.Errorf("thrift server init error[%v]",err)
	//	return
	//}

	server=&Server{
		httpServer:httpServer,
		thriftServer:thriftServer,
	}
	return
}

func (this *Server) Stop() {
	log.Info("Server stop")

	if this.httpServer !=nil {
		this.httpServer.Stop()
	}

	if this.thriftServer !=nil {
		this.thriftServer.Stop()
	}
}
