package defines

type AllocStrategy interface {
	Write(key int64,cookie int32,img []byte) (string,error)
	Read(url string) ([]byte,error)
	Delete(url string) error
	Stop()
}