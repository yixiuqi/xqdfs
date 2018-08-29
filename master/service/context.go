package service

import (
	"xqdfs/master/conf"
)

type Context struct {
	Conf *conf.Config
}

func NewContext(c *conf.Config) *Context {
	context:=new(Context)
	context.Conf=c
	return context
}
