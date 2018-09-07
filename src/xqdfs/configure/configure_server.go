package configure

import (
	"xqdfs/configure/ssdb"
	"xqdfs/configure/defines"
)

type ConfigureServer struct {
	configure defines.Configure
}

func NewConfigureServer(param string) (*ConfigureServer,error) {
	s:=&ConfigureServer{
		configure:ssdb.NewConfigureSSDB(param),
	}
	return s,nil
}

func (this *ConfigureServer) ParamGet(key string) (string,error) {
	return this.configure.ParamGet(key)
}

func (this *ConfigureServer) ParamSet(key string,value string) error {
	return this.configure.ParamSet(key,value)
}

func (this *ConfigureServer) ParamSetx(key string,value string,ttl int) error {
	return this.configure.ParamSetx(key,value,ttl)
}

func (this *ConfigureServer) StorageAdd(s *defines.StorageDal) error {
	return this.configure.StorageAdd(s)
}

func (this *ConfigureServer) StorageRemove(sid int32) error {
	return this.configure.StorageRemove(sid)
}

func (this *ConfigureServer) StorageGet(sid int32) (*defines.StorageDal,error) {
	return this.configure.StorageGet(sid)
}

func (this *ConfigureServer) StorageGetAll() ([]*defines.StorageDal,error) {
	return this.configure.StorageGetAll()
}

func (this *ConfigureServer) GroupAdd(g *defines.GroupDal) error {
	return this.configure.GroupAdd(g)
}

func (this *ConfigureServer) GroupRemove(gid int32) error {
	return this.configure.GroupRemove(gid)
}

func (this *ConfigureServer) GroupEdit(g *defines.GroupDal) (err error) {
	return this.configure.GroupEdit(g)
}

func (this *ConfigureServer) GroupGet(gid int32) (g *defines.GroupDal,err error) {
	return this.configure.GroupGet(gid)
}

func (this *ConfigureServer) GroupGetAll() (g []*defines.GroupDal,err error) {
	return this.configure.GroupGetAll()
}

func (this *ConfigureServer) Stop() {
	this.configure.Stop()
}