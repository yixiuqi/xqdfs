package thrift

import (
	"net"
	"sync"
	"time"
	"strconv"

	"xqdfs/utils/log"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type ThriftClientItem struct{
	id int
	socket *thrift.TSocket
	client *ThriftServiceClient
	host string
	port int
	timeout time.Duration
	lock sync.Mutex
}

func NewThriftClientItem(host string,port int,id int) *ThriftClientItem{
	p:=new(ThriftClientItem)
	p.init(host,port,id)
	return p
}

func (this *ThriftClientItem) init(host string,port int,id int){
	this.id=id
	this.host=host
	this.port=port
	this.timeout=30*time.Second

	if this.socket!=nil{
		this.socket.Close()
		this.socket=nil
	}

	sock,err:=thrift.NewTSocketTimeout(net.JoinHostPort(this.host, strconv.Itoa(this.port)),this.timeout)
	if err != nil {
		log.Error("thrift连接服务端失败:"+err.Error())
	}else{
		this.socket=sock
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	this.client = NewThriftServiceClientFactory(this.socket, protocolFactory)
	if err := this.socket.Open(); err != nil {
		log.Error("thrift连接服务端失败:"+err.Error())
	}
}

func (this *ThriftClientItem) destroy(){
	if this.socket!=nil{
		this.socket.Close()
		this.socket=nil
	}
}

func (this *ThriftClientItem) Send(cmd string, body *ThriftCmdStruct) (*ThriftResultStrcut,error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	ret, err := this.client.DoProcess(cmd, "XQ_20100927", body)
	if err == nil {
		return ret,nil
	} else {
		this.init(this.host, this.port, this.id)
		return nil,err
	}
}
