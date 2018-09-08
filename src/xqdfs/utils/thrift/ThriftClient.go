package thrift

import (
	"sync"
)

type ThriftClient struct {
	host string
	port int
	link []*ThriftClientItem
	id   int
	lock sync.Mutex
}

/*
@brief
	新建一个thrift客户端
@param
	host-服务器地址
	port-服务器端口
	count-连接池数目
*/
func NewThriftClient(host string, port int, count int) *ThriftClient {
	client := new(ThriftClient)
	client.host = host
	client.port = port
	client.link = make([]*ThriftClientItem, count)
	for i := 0; i < count; i++ {
		client.link[i] = NewThriftClientItem(client.host, client.port, i+1)
	}
	return client
}

func (this *ThriftClient) Destroy() {
	for i := 0; i < len(this.link); i++ {
		this.link[i].destroy()
	}
}

func (this *ThriftClient) Send(cmd string, body *ThriftCmdStruct) (*ThriftResultStrcut,error) {
	var link *ThriftClientItem

	this.lock.Lock()
	this.id++
	if this.id >= len(this.link) {
		this.id = 0
	}
	link = this.link[this.id]
	this.lock.Unlock()

	return link.Send(cmd, body)
}
