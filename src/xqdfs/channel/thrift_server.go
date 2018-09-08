package channel

import (
	"fmt"

	"xqdfs/utils/log"
	thriftdef "xqdfs/utils/thrift"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type ThriftServer struct {
}

func NewThriftServer(port int) (*ThriftServer,error) {
	server:=&ThriftServer{}
	go server.process(port)
	return server,nil
}

func (this *ThriftServer) process(port int) {
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTBinaryProtocolFactory(false, false)
	transport, _ := thrift.NewTServerSocket(fmt.Sprintf(":%d",port))
	processor := thriftdef.NewThriftServiceProcessor(this)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	server.Serve()
}

func (this *ThriftServer) DoProcess(cmd string, seq string, body *thriftdef.ThriftCmdStruct) (r *thriftdef.ThriftResultStrcut, err error) {
	return nil,nil
}

func (this *ThriftServer) Stop() {
	log.Info("ThriftServer stop")
}
