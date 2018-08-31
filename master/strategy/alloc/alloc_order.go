package order

import (
	"fmt"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/configure"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/master/conf"
	"xqdfs/master/strategy/tool"
	"xqdfs/master/strategy/defines"
)

type AllocOrder struct {
	conf *conf.Config
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	selectWritableVolume *SelectWritableVolume
}

func NewAllocOrder(conf *conf.Config,configureServer *configure.ConfigureServer,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) *AllocOrder {
	s:=&AllocOrder{
		conf:conf,
		configureServer:configureServer,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		selectWritableVolume:NewSelectWritableVolume(discoveryServer),
	}
	return s
}

func (this *AllocOrder) Write(key int64,cookie int32,img []byte) (string,error) {
	removeVolumes:=make([]*defines.WritableVolume,0)
	volume,err:=this.selectWritableVolume.SelectWritableVolume(this.conf,int32(len(img)),removeVolumes)
	if err!=nil {
		log.Debug(err)
		return "",err
	}

	vid:=volume.VolumeId
	host:=volume.StorageAddr
	err=this.proxyStorage.Upload(host,vid,key,cookie,img,true)

	//try three times
	count:=0
	for err==errors.ErrSuperBlockNoSpace&&count<this.conf.AllocStrategy.OrderConsumeCount*2 {
		removeVolumes=append(removeVolumes,volume)
		count++
		log.Debugf("[%s][%v] ErrSuperBlockNoSpace try[%d][%v] [%v]",host,vid,count,volume,removeVolumes)
		volume,err=this.selectWritableVolume.SelectWritableVolume(this.conf,int32(len(img)),removeVolumes)
		if err!=nil {
			log.Debug(err)
			return "",err
		}
		vid=volume.VolumeId
		host=volume.StorageAddr
		err=this.proxyStorage.Upload(host,vid,key,cookie,img,true)
	}

	if err!=nil{
		log.Errorf("[%s][%v] error:[%v]",host,vid,err)
		return "",err
	}else{
		id:=fmt.Sprintf("%d,%d,%d,%d,%d",volume.GroupId,volume.StorageId,volume.VolumeId,key,cookie)
		return id,nil
	}
}

func (this *AllocOrder) Read(url string) ([]byte,error) {
	groups:=this.discoveryServer.Groups()
	if groups==nil {
		log.Debug("groups is null")
		return nil,errors.ErrGroupIsEmpty
	}

	location,err:=tool.UrlToLocation(url)
	if err!=nil {
		log.Error(err)
		return nil,errors.ErrParameterError
	}

	source,err:=selectImageSource(groups,location)
	if err!=nil{
		log.Debug(err)
		return nil,err
	}else{
		log.Debugf("image source group[%d] storage[%d][%s] volume[%d] key[%d] cookie[%d]",
			location.GroupId,location.StorageId,source.Host,location.VolumeId,location.Key,location.Cookie)
	}

	img,err:=this.proxyStorage.Get(source.Host,location.VolumeId,location.Key,location.Cookie)
	if err!=nil {
		log.Error(err)
		return nil,err
	}else{
		return img,nil
	}
}

func (this *AllocOrder) Delete(url string) error {
	groups:=this.discoveryServer.Groups()
	if groups==nil {
		log.Debug("groups is null")
		return errors.ErrGroupIsEmpty
	}

	location,err:=tool.UrlToLocation(url)
	if err!=nil {
		log.Error(err)
		return errors.ErrParameterError
	}

	source,err:=selectImageSource(groups,location)
	if err!=nil{
		log.Debug(err)
		return err
	}else{
		log.Debugf("image source group[%d] storage[%d][%s] volume[%d] key[%d] cookie[%d]",
			location.GroupId,location.StorageId,source.Host,location.VolumeId,location.Key,location.Cookie)
	}

	err=this.proxyStorage.Delete(source.Host,location.VolumeId,location.Key,true)
	if err!=nil {
		log.Error(err)
		return err
	}else{
		return nil
	}
}

func (this *AllocOrder) Stop() {
	log.Info("AllocOrder stop")
}
