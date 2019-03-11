package defines

import "github.com/Jeffail/gabs"

type Proxy interface {
	//上传由于图片比较大，解码耗时，所以后面直接传递gabs结构给后续处理
	//body是包含了img的json
	Upload(host string,body *gabs.Container,vid int32,key int64,cookie int32,img []byte,replication bool) error
	Get(host string,vid int32,key int64,cookie int32) ([]byte,error)
	Delete(host string,vid int32,key int64,replication bool) error
	StorageInit(host string,replication bool) error
	StorageVolumeCompact(host string,vid int32,replication bool) error
	StorageVolumeClear(host string,vid int32,replication bool) error
	StorageGetConfigure(host string) (*gabs.Container,error)
	Stop()
}
