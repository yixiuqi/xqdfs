package service

import (
	"xqdfs/storage/store"
	"xqdfs/storage/conf"
	"xqdfs/storage/replication"
)

type Context struct {
	Conf *conf.Config
	Store *store.Store
	ReplicationServer *replication.ReplicationServer
}

func NewContext(c *conf.Config,s *store.Store,replicationServer *replication.ReplicationServer) *Context {
	context:=new(Context)
	context.Conf=c
	context.Store=s
	context.ReplicationServer=replicationServer
	return context
}
