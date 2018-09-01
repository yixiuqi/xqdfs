package order

import (
	"fmt"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/discovery"
	"xqdfs/configure"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/utils/plugin"
	"xqdfs/master/strategy/tool"
	"xqdfs/master/strategy/defines"
)

const(
	AllocOrderMinFreeSpace 	= 	"AllocOrderMinFreeSpace"			// 每个卷预留多少空间 default:104857600
	AllocOrderConsumeCount 	= 	"AllocOrderConsumeCount"			// 同时default:3
)

type AllocOrder struct {
	configureServer *configure.ConfigureServer
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	selectWritableVolume *SelectWritableVolume
	orderMinFreeSpace int64
	orderConsumeCount int
}

func NewAllocOrder() (*AllocOrder,error) {
	var conf *configure.ConfigureServer
	if s:=plugin.PluginGetObject(plugin.PluginConfigure);s==nil {
		log.Errorf("%s no support",plugin.PluginConfigure)
		return nil,errors.ErrNoSupport
	}else{
		conf=s.(*configure.ConfigureServer)
	}

	var discoveryServer *discovery.DiscoveryServer
	if d:=plugin.PluginGetObject(plugin.PluginDiscoveryServer);d==nil {
		log.Errorf("%s no support",plugin.PluginDiscoveryServer)
		return nil,errors.ErrNoSupport
	}else{
		discoveryServer=d.(*discovery.DiscoveryServer)
	}

	var proxyStorage *proxy.ProxyStorage
	if p:=plugin.PluginGetObject(plugin.PluginProxyStorage);p==nil {
		log.Errorf("%s no support",plugin.PluginProxyStorage)
		return nil,errors.ErrNoSupport
	}else{
		proxyStorage=p.(*proxy.ProxyStorage)
	}

	var orderMinFreeSpace int64=104857600
	value,err:=conf.ParamGet(AllocOrderMinFreeSpace)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=conf.ParamSet(AllocOrderMinFreeSpace,"104857600")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		orderMinFreeSpace,err=helper.StringToInt64(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	orderConsumeCount:=3
	value,err=conf.ParamGet(AllocOrderConsumeCount)
	if err!=nil{
		if err==errors.ErrParamNotExist{
			err=conf.ParamSet(AllocOrderConsumeCount,"3")
			if err!=nil{
				return nil,err
			}
		}else{
			log.Error(err)
			return nil,err
		}
	}else{
		orderConsumeCount,err=helper.StringToInt(value)
		if err!=nil{
			log.Error(err)
			return nil,err
		}
	}

	log.Infof("%s[%d]",AllocOrderMinFreeSpace,orderMinFreeSpace)
	log.Infof("%s[%d]",AllocOrderConsumeCount,orderConsumeCount)
	s:=&AllocOrder{
		configureServer:conf,
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		selectWritableVolume:NewSelectWritableVolume(discoveryServer),
		orderMinFreeSpace:orderMinFreeSpace,
		orderConsumeCount:orderConsumeCount,
	}
	ServiceAllocOrderSetup(s)
	return s,nil
}

func (this *AllocOrder) Write(key int64,cookie int32,img []byte) (string,error) {
	removeVolumes:=make([]*defines.WritableVolume,0)
	volume,err:=this.selectWritableVolume.SelectWritableVolume(this.orderMinFreeSpace,this.orderConsumeCount,int32(len(img)),removeVolumes)
	if err!=nil {
		log.Debug(err)
		return "",err
	}

	vid:=volume.VolumeId
	host:=volume.StorageAddr
	err=this.proxyStorage.Upload(host,vid,key,cookie,img,true)

	//try three times
	count:=0
	for err==errors.ErrSuperBlockNoSpace&&count<this.orderConsumeCount*2 {
		removeVolumes=append(removeVolumes,volume)
		count++
		log.Debugf("[%s][%v] ErrSuperBlockNoSpace try[%d][%v] [%v]",host,vid,count,volume,removeVolumes)
		volume,err=this.selectWritableVolume.SelectWritableVolume(this.orderMinFreeSpace,this.orderConsumeCount,int32(len(img)),removeVolumes)
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

func (this *AllocOrder) AllocOrderMinFreeSpaceGet() int64 {
	return this.orderMinFreeSpace
}

func (this *AllocOrder) AllocOrderMinFreeSpaceSet(orderMinFreeSpace int64) error {
	err:=this.configureServer.ParamSet(AllocOrderMinFreeSpace,fmt.Sprintf("%d",orderMinFreeSpace))
	if err!=nil{
		log.Error(err)
		return err
	}else{
		this.orderMinFreeSpace=orderMinFreeSpace
		return nil
	}
}

func (this *AllocOrder) AllocOrderConsumeCountGet() int {
	return this.orderConsumeCount
}

func (this *AllocOrder) AllocOrderConsumeCountSet(orderConsumeCount int) error {
	err:=this.configureServer.ParamSet(AllocOrderConsumeCount,fmt.Sprintf("%d",orderConsumeCount))
	if err!=nil{
		log.Error(err)
		return err
	}else{
		this.orderConsumeCount=orderConsumeCount
		return nil
	}
}
