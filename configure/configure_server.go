package configure

import (
	"xqdfs/configure/defines"
	"xqdfs/configure/ssdb"
)

type ConfigureServer struct {
	configure defines.GroupConf
}

func NewConfigureServer(param string) (*ConfigureServer,error) {
	s:=&ConfigureServer{
		configure:ssdb.NewGroupConfSSDB(param),
	}
	return s,nil
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