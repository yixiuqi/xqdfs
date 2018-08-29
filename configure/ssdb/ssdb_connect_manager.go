package ssdb

import (
	"strings"
	"errors"
	"fmt"
	"net"
	"sync/atomic"

	"xqdfs/utils/log"
	"xqdfs/utils/helper"
)

type SSDBConnectMgr struct {
	pos uint64
	addr []string
	port []int
}

var gSSDBConnectMgr *SSDBConnectMgr
func SSDBConnectMgrInstance() *SSDBConnectMgr {
	if gSSDBConnectMgr == nil {
		gSSDBConnectMgr = new(SSDBConnectMgr)
		log.Info(helper.NewStringBuilder().Append("create SSDBConnectMgr").ToString())
	}

	return gSSDBConnectMgr
}

func (this *SSDBConnectMgr) Init(addr string) error{
	items:=strings.Split(addr,",")
	if items==nil||len(items)<=0{
		return errors.New("SSDB地址配置错误")
	}

	this.addr=make([]string,len(items))
	this.port=make([]int,len(items))

	for i:=0;i<len(items);i++{
		kv:=strings.Split(items[i],":")
		if kv==nil||len(kv)!=2{
			return errors.New("SSDB地址配置错误")
		}

		this.addr[i]=kv[0]
		port,err:=helper.StringToInt(kv[1])
		if err!=nil{
			return errors.New("SSDB地址配置错误")
		}

		this.port[i]=port
	}

	return nil
}

func (this *SSDBConnectMgr) Dump() {
	if this.addr==nil||this.port==nil{
		fmt.Println("null")
		return
	}

	for i:=0;i<len(this.addr);i++{
		fmt.Println(this.addr[i],":",this.port[i])
	}
}

func (this *SSDBConnectMgr) getConnect() (*SSDBConnect, error) {
	len:=len(this.addr)

	pos:=this.pos%uint64(len)
	conn,err:=this.connect(this.addr[pos],this.port[pos])
	if err!=nil{
		atomic.AddUint64(&this.pos,1)
		pos:=this.pos%uint64(len)
		conn,err=this.connect(this.addr[pos],this.port[pos])
		if err!=nil{
			return nil,err
		}else{
			debug:=helper.NewStringBuilder().Append("获取连接").Append(this.addr[pos]).Append(":").Append(this.port[pos]).ToString()
			log.Debug(debug)
			return conn,err
		}
	}else{
		debug:=helper.NewStringBuilder().Append("获取连接").Append(this.addr[pos]).Append(":").Append(this.port[pos]).ToString()
		log.Trace(debug)
		return conn,err
	}
}

func (this *SSDBConnectMgr) connect(host string,port int) (*SSDBConnect, error) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}
	sock, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	conn:=NewSSDBConnect(host,port,sock)
	return conn, nil
}