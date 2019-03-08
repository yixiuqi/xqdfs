package configure

import (
	"encoding/json"

	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/configure/ssdb"
	"xqdfs/configure/defines"
)

const(
	HashXQDfsGroup 	= "HashXQDfsGroup"
	HashXQDfsStorage 	= "HashXQDfsStorage"
	HashXQDfsConfig 	= "HashXQDfsConfig"
)

type ConfigureServer struct {
	kv *ssdb.SSDBKV
	hset *ssdb.SSDBHash
}

//param is remote config server's connection param
func NewConfigureServer(param string) (*ConfigureServer,error) {
	log.Info("NewConfigureServer param is ",param)
	connMgr:=&ssdb.SSDBConnectMgr{}
	err:=connMgr.Init(param)
	if err!=nil {
		log.Error(err)
		return nil,err
	}

	server:=&ConfigureServer{
		kv:ssdb.NewSSDBKV(connMgr),
		hset:ssdb.NewSSDBHash(connMgr),
	}
	return server,nil
}

func (this *ConfigureServer) KVGet(key string) (string,error) {
	return this.kv.Get(key)
}

func (this *ConfigureServer) KVSet(key string,value string) error {
	return this.kv.Set(key,value)
}

func (this *ConfigureServer) KVSetx(key string,value string,ttl int) error {
	return this.kv.Setx(key,value,ttl)
}

func (this *ConfigureServer) ConfigGet(key string) (string,error) {
	return this.hset.HGet(HashXQDfsConfig,key)
}

func (this *ConfigureServer) ConfigSet(key string,value string) error {
	return this.hset.HSet(HashXQDfsConfig,key,value)
}

func (this *ConfigureServer) StorageAdd(s *defines.StorageDal) error {
	param,err:=json.Marshal(s)
	if err!=nil {
		log.Warn(err)
		return err
	}

	return this.hset.HSet(HashXQDfsStorage,helper.Int32ToString(s.Id),string(param))
}

func (this *ConfigureServer) StorageRemove(sid int32) error {
	return this.hset.HDel(HashXQDfsStorage,helper.Int32ToString(sid))
}

func (this *ConfigureServer) StorageGet(sid int32) (*defines.StorageDal,error) {
	param,err:=this.hset.HGet(HashXQDfsStorage,helper.Int32ToString(sid))
	if err!=nil {
		log.Warn(err)
		return nil,err
	}
	if param=="" {
		return nil,nil
	}

	s:=defines.NewStorageDal()
	err=json.Unmarshal([]byte(param),s)
	if err!=nil {
		log.Warn(err)
		return nil,err
	}

	return s,nil
}

func (this *ConfigureServer) StorageGetAll() ([]*defines.StorageDal,error) {
	items,err:=this.hset.HGetAll(HashXQDfsStorage)
	if err!=nil {
		log.Warn(err)
		return nil,err
	}

	if items==nil{
		return nil,nil
	}

	s:=make([]*defines.StorageDal,0)
	for _,v:=range items {
		one:=defines.NewStorageDal()
		e:=json.Unmarshal([]byte(v.Value),one)
		if e==nil{
			s=append(s,one)
		}else{
			log.Warn(err)
		}
	}

	return s,nil
}

func (this *ConfigureServer) GroupAdd(g *defines.GroupDal) error {
	param,err:=json.Marshal(g)
	if err!=nil {
		log.Warn(err)
		return err
	}

	return this.hset.HSet(HashXQDfsGroup,helper.Int32ToString(g.Id),string(param))
}

func (this *ConfigureServer) GroupRemove(gid int32) error {
	return this.hset.HDel(HashXQDfsGroup,helper.Int32ToString(gid))
}

func (this *ConfigureServer) GroupEdit(g *defines.GroupDal) error {
	param,err:=json.Marshal(g)
	if err!=nil {
		log.Warn(err)
		return err
	}

	return this.hset.HSet(HashXQDfsGroup,helper.Int32ToString(g.Id),string(param))
}

func (this *ConfigureServer) GroupGet(gid int32) (*defines.GroupDal,error) {
	param,err:=this.hset.HGet(HashXQDfsGroup,helper.Int32ToString(gid))
	if err!=nil {
		log.Warn(err)
		return nil,err
	}

	if param==""{
		return nil,nil
	}

	g:=defines.NewGroupDal()
	err=json.Unmarshal([]byte(param),g)
	if err!=nil {
		log.Warn(err)
		return nil,err
	}

	return g,nil
}

func (this *ConfigureServer) GroupGetAll() ([]*defines.GroupDal,error) {
	items,err:=this.hset.HGetAll(HashXQDfsGroup)
	if err!=nil {
		log.Warn(err)
		return nil,err
	}
	if items==nil{
		return nil,nil
	}

	g:=make([]*defines.GroupDal,0)
	for _,v:=range items {
		one:=defines.NewGroupDal()
		e:=json.Unmarshal([]byte(v.Value),one)
		if e==nil{
			g=append(g,one)
		}else{
			log.Warn(err)
		}
	}

	return g,nil
}

func (this *ConfigureServer) Stop() {
	log.Info("ConfigureServer stop")
}