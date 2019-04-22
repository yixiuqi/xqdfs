package alloc

import (
	"fmt"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/discovery"
	"xqdfs/utils/log"
	"xqdfs/utils/plugin"
	"xqdfs/master/configure"
	"xqdfs/master/strategy/tool"
	"xqdfs/master/strategy/defines"

	"github.com/Jeffail/gabs"
)

type AllocOrder struct {
	discoveryServer *discovery.DiscoveryServer
	proxyStorage *proxy.ProxyStorage
	selectWritableVolume *SelectWritableVolume

	uploadErrorProcess *UploadErrorProcess
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

	s:=&AllocOrder{
		discoveryServer:discoveryServer,
		proxyStorage:proxyStorage,
		selectWritableVolume:NewSelectWritableVolume(discoveryServer),
		uploadErrorProcess:NewUploadErrorProcess(proxyStorage),
	}

	setupConfigureSvr(conf)
	setupService()
	return s,nil
}

func (this *AllocOrder) Write(key int64,cookie int32,body *gabs.Container) (string,error) {
	//大致算一下jason里面图片实际长度
	imgLen:=int32(len(body.Bytes())*3/4)
	removeVolumes:=make([]*defines.WritableVolume,0)
	volume,err:=this.selectWritableVolume.SelectWritableVolume(imgLen,removeVolumes)
	if err!=nil {
		log.Debug(err)
		return "",err
	}

	vid:=volume.VolumeId
	host:=volume.StorageAddr
	err=this.proxyStorage.Upload(host,body,vid,key,cookie,nil,true)
	if err==errors.ErrRpc {
		log.Warnf("Upload[%s] error",host)
		this.uploadErrorProcess.RollBack(host,vid,key)
	}

	//try 3 times
	count:=0
	for err==errors.ErrSuperBlockNoSpace&&count< 3 {
		removeVolumes=append(removeVolumes,volume)
		count++
		log.Debugf("[%s][%v] ErrSuperBlockNoSpace try[%d][%v] [%v]",host,vid,count,volume,removeVolumes)
		volume,err=this.selectWritableVolume.SelectWritableVolume(imgLen,removeVolumes)
		if err!=nil {
			log.Debug(err)
			return "",err
		}
		vid=volume.VolumeId
		host=volume.StorageAddr
		err=this.proxyStorage.Upload(host,body,vid,key,cookie,nil,true)
		if err==errors.ErrRpc {
			log.Warnf("Upload[%s] error",host)
			this.uploadErrorProcess.RollBack(host,vid,key)
		}
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
	log.Info("AllocOrder stop->")
	this.uploadErrorProcess.Stop()
	log.Info("AllocOrder stop-<")
}