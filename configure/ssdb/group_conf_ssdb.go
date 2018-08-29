package ssdb

import (
	"encoding/json"

	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/configure/defines"
)

const(
	HashNameGroup 	= "xqdfs_group"
	HashNameStorage 	= "xqdfs_storage"
)

type GroupConfSSDB struct {
}

func NewGroupConfSSDB(addr string) (g *GroupConfSSDB) {
	g=new(GroupConfSSDB)
	SSDBConnectMgrInstance().Init(addr)
	return
}

func (this *GroupConfSSDB) StorageAdd(s *defines.StorageDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(s)
	if err==nil{
		err=hset.HSet(HashNameStorage,helper.Int32ToString(s.Id),string(param))
	}

	return
}

func (this *GroupConfSSDB) StorageRemove(sid int32) (err error) {
	hset:=NewSSDBHash()
	err=hset.HDel(HashNameStorage,helper.Int32ToString(sid))
	return
}

func (this *GroupConfSSDB) StorageGet(sid int32) (s *defines.StorageDal,err error) {
	hset:=NewSSDBHash()
	var param string
	param,err=hset.HGet(HashNameStorage,helper.Int32ToString(sid))
	if err!=nil||param==""{
		return
	}

	s=defines.NewStorageDal()
	err=json.Unmarshal([]byte(param),s)
	return
}

func (this *GroupConfSSDB) StorageGetAll() (s []*defines.StorageDal,err error) {
	hset:=NewSSDBHash()
	var items []SSDBHashItem
	items,err=hset.HGetAll(HashNameStorage)
	if err!=nil||items==nil{
		return
	}

	s=make([]*defines.StorageDal,0)
	for _,v:=range items {
		one:=defines.NewStorageDal()
		e:=json.Unmarshal([]byte(v.Value),one)
		if e==nil{
			s=append(s,one)
		}else{
			log.Error(err)
		}
	}

	return
}

func (this *GroupConfSSDB) GroupAdd(g *defines.GroupDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(g)
	if err==nil{
		err=hset.HSet(HashNameGroup,helper.Int32ToString(g.Id),string(param))
	}

	return
}

func (this *GroupConfSSDB) GroupRemove(gid int32) (err error){
	hset:=NewSSDBHash()
	err=hset.HDel(HashNameGroup,helper.Int32ToString(gid))
	return
}

func (this *GroupConfSSDB) GroupEdit(g *defines.GroupDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(g)
	if err==nil{
		err=hset.HSet(HashNameGroup,helper.Int32ToString(g.Id),string(param))
	}

	return
}

func (this *GroupConfSSDB) GroupGet(gid int32) (g *defines.GroupDal,err error) {
	hset:=NewSSDBHash()
	var param string
	param,err=hset.HGet(HashNameGroup,helper.Int32ToString(gid))
	if err!=nil||param==""{
		return
	}

	g=defines.NewGroupDal()
	err=json.Unmarshal([]byte(param),g)
	return
}

func (this *GroupConfSSDB) GroupGetAll() (g []*defines.GroupDal,err error) {
	hset:=NewSSDBHash()
	var items []SSDBHashItem
	items,err=hset.HGetAll(HashNameGroup)
	if err!=nil||items==nil{
		return
	}

	g=make([]*defines.GroupDal,0)
	for _,v:=range items {
		one:=defines.NewGroupDal()
		e:=json.Unmarshal([]byte(v.Value),one)
		if e==nil{
			g=append(g,one)
		}else{
			log.Error(err)
		}
	}

	return
}

func (this *GroupConfSSDB) Stop() {
	log.Info("GroupConfSSDB stop")
}