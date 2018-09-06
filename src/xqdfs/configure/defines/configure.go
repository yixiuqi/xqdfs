package defines

const(
	HashNameGroup 	= "xqdfs_group"
	HashNameStorage 	= "xqdfs_storage"
)

type Configure interface {
	ParamGet(key string) (string,error)
	ParamSet(key string,value string) error
	StorageAdd(s *StorageDal) error
	StorageRemove(sid int32) error
	StorageGet(sid int32) (*StorageDal,error)
	StorageGetAll() ([]*StorageDal,error)
	GroupAdd(g *GroupDal) error
	GroupRemove(gid int32) error
	GroupEdit(g *GroupDal) error
	GroupGet(gid int32) (*GroupDal,error)
	GroupGetAll() ([]*GroupDal,error)
	Stop()
}