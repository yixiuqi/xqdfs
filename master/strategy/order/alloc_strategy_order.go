package order

/*
* @brief 所有组、Storage、卷按Id排序，只要空间足够就写入
* @author    yimin
* @email     yimin@cloudwalk.cn
* @date      2018/08/06
*/

import (
	"fmt"

	"xqdfs/master/conf"
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/master/strategy/tool"
	"xqdfs/master/strategy/defines"
	"xqdfs/discovery"
	"xqdfs/proxy"
)

type AllocStrategyOrder struct {
	conf *conf.Config
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	orderClearTask *OrderClearTask
	selectWritableVolume *SelectWritableVolume
}

func NewAllocStrategyOrder(conf *conf.Config,discoveryServer *discovery.DiscoveryServer,proxyStorage *proxy.ProxyStorage) *AllocStrategyOrder {
	s:=&AllocStrategyOrder{
		conf:conf,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		orderClearTask:NewOrderClearTask(conf,discoveryServer,proxyStorage),
		selectWritableVolume:NewSelectWritableVolume(discoveryServer),
	}
	return s
}

func (this *AllocStrategyOrder) Write(key int64,cookie int32,img []byte) (string,error) {
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

func (this *AllocStrategyOrder) Read(url string) ([]byte,error) {
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

func (this *AllocStrategyOrder) Delete(url string) error {
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

func (this *AllocStrategyOrder) Stop() {
	log.Info("AllocStrategyOrder stop")
	this.orderClearTask.Stop()
}
