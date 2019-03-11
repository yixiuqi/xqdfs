package defines

import "github.com/Jeffail/gabs"

type AllocStrategy interface {
	Write(key int64,cookie int32,body *gabs.Container) (string,error)
	Read(url string) ([]byte,error)
	Delete(url string) error
	Stop()
}