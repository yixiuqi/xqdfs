package ssdb

import (
	"encoding/json"

	"xqdfs/utils/helper"
	"xqdfs/utils/log"
	"xqdfs/configure/defines"
)

type ConfigureSSDB struct {
}

func NewConfigureSSDB(addr string) *ConfigureSSDB {
	SSDBConnectMgrInstance().Init(addr)
	return &ConfigureSSDB{}
}

func (this *ConfigureSSDB) ParamGet(key string) (string,error) {
	kv:=NewSSDBKV()
	return kv.Get(key)
}

func (this *ConfigureSSDB) ParamSet(key string,value string) error {
	kv:=NewSSDBKV()
	return kv.Set(key,value)
}

func (this *ConfigureSSDB) ParamSetx(key string,value string,ttl int) error {
	kv:=NewSSDBKV()
	return kv.Setx(key,value,ttl)
}

func (this *ConfigureSSDB) StorageAdd(s *defines.StorageDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(s)
	if err==nil{
		err=hset.HSet(defines.HashNameStorage,helper.Int32ToString(s.Id),string(param))
	}

	return
}

func (this *ConfigureSSDB) StorageRemove(sid int32) (err error) {
	hset:=NewSSDBHash()
	err=hset.HDel(defines.HashNameStorage,helper.Int32ToString(sid))
	return
}

func (this *ConfigureSSDB) StorageGet(sid int32) (s *defines.StorageDal,err error) {
	hset:=NewSSDBHash()
	var param string
	param,err=hset.HGet(defines.HashNameStorage,helper.Int32ToString(sid))
	if err!=nil||param==""{
		return
	}

	s=defines.NewStorageDal()
	err=json.Unmarshal([]byte(param),s)
	return
}

func (this *ConfigureSSDB) StorageGetAll() (s []*defines.StorageDal,err error) {
	hset:=NewSSDBHash()
	var items []SSDBHashItem
	items,err=hset.HGetAll(defines.HashNameStorage)
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

func (this *ConfigureSSDB) GroupAdd(g *defines.GroupDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(g)
	if err==nil{
		err=hset.HSet(defines.HashNameGroup,helper.Int32ToString(g.Id),string(param))
	}

	return
}

func (this *ConfigureSSDB) GroupRemove(gid int32) (err error){
	hset:=NewSSDBHash()
	err=hset.HDel(defines.HashNameGroup,helper.Int32ToString(gid))
	return
}

func (this *ConfigureSSDB) GroupEdit(g *defines.GroupDal) (err error) {
	hset:=NewSSDBHash()
	var param []byte
	param,err=json.Marshal(g)
	if err==nil{
		err=hset.HSet(defines.HashNameGroup,helper.Int32ToString(g.Id),string(param))
	}

	return
}

func (this *ConfigureSSDB) GroupGet(gid int32) (g *defines.GroupDal,err error) {
	hset:=NewSSDBHash()
	var param string
	param,err=hset.HGet(defines.HashNameGroup,helper.Int32ToString(gid))
	if err!=nil||param==""{
		return
	}

	g=defines.NewGroupDal()
	err=json.Unmarshal([]byte(param),g)
	return
}

func (this *ConfigureSSDB) GroupGetAll() (g []*defines.GroupDal,err error) {
	hset:=NewSSDBHash()
	var items []SSDBHashItem
	items,err=hset.HGetAll(defines.HashNameGroup)
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

func (this *ConfigureSSDB) Stop() {
	log.Info("ConfigureSSDB stop")
}