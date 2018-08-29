package service

import (
	"xqdfs/storage/store"
	"xqdfs/storage/conf"
)

type Context struct {
	Conf *conf.Config
	Store *store.Store
}

func NewContext(c *conf.Config,s *store.Store) *Context {
	context:=new(Context)
	context.Conf=c
	context.Store=s
	return context
}
