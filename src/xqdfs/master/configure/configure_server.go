package configure

import (
	"time"
	"encoding/json"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/master/configure/ssdb"
	"xqdfs/master/configure/defines"

	"github.com/json-iterator/go"
)

const(
	HashXQDfsGroup 	= "HashXQDfsGroup"
	HashXQDfsStorage 	= "HashXQDfsStorage"
	HashXQDfsConfig 	= "HashXQDfsConfig"
)

type ConfigureServer struct {
	kv *ssdb.SSDBKV
	hset *ssdb.SSDBHash
	signal chan int
	change chan int
	params map[string]string
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
		signal:make(chan int),
		change:make(chan int),
		params:make(map[string]string,0),
	}
	server.init()
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
	value,ok:=this.params[key]
	if ok {
		return value,nil
	}else{
		return "",errors.ErrParamNotExist
	}
}

func (this *ConfigureServer) ConfigSet(key string,value string) error {
	err:=this.hset.HSet(HashXQDfsConfig,key,value)
	if err!=nil {
		log.Warn(err)
		return err
	}

	defer helper.HandleErr()
	this.change<-1
	return nil
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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
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
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		e:=json.Unmarshal([]byte(v.Value),one)
		if e==nil{
			g=append(g,one)
		}else{
			log.Warn(err)
		}
	}

	return g,nil
}

func (this *ConfigureServer) getConfigreParam() {
	all,err:=this.hset.HGetAll(HashXQDfsConfig)
	if err!=nil {
		log.Warn(err)
		this.params=make(map[string]string,0)
		return
	}
	params:=make(map[string]string,0)
	for _,v:=range all {
		params[v.Key]=v.Value
	}
	this.params=params
}

func (this *ConfigureServer) init() {
	this.getConfigreParam()
	log.Info("------ Configure params ------")
	for k,v:=range this.params {
		log.Infof("%s:%s",k,v)
	}
	go func(){
		for {
			select {
			case <-time.After(time.Second * 10):
			case <-this.signal:
				this.signal<-1
				return
			case <-this.change:
				log.Debug("change params")
			}

			this.getConfigreParam()
		}
	}()
}

func (this *ConfigureServer) Stop() {
	log.Info("ConfigureServer stop->")
	this.signal<-1
	<-this.signal
	close(this.signal)
	close(this.change)
	log.Info("ConfigureServer stop-<")
}