package defines

type GroupConf interface {
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