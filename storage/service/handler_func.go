package service

type HandlerFunc func(context *Context,m map[string]interface{}) interface{}
