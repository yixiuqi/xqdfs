package service

import (

)

type HandlerFunc func(context *Context,m map[string]interface{}) interface{}

const(
	Success 	= 	0
	Failed 	=	1
)
