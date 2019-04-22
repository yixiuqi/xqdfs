package http

import (
	"encoding/base64"

	"xqdfs/errors"
	"xqdfs/constant"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"

	"github.com/Jeffail/gabs"
)

type ProxyHttp struct {
}

func NewProxyHttp() *ProxyHttp {
	p:=new(ProxyHttp)
	return p
}

func (this *ProxyHttp) Upload(host string,body *gabs.Container,vid int32,key int64,cookie int32,img []byte,replication bool) error {
	if body==nil {
		body=gabs.New()
	}

	body.Set(helper.UUIDBuild(),"seq")
	body.Set(vid,"vid")
	body.Set(key,"key")
	body.Set(cookie,"cookie")
	if img!=nil {
		body.Set(img,"img")
	}
	body.Set(replication,"replication")
	url:="http://"+host+constant.CmdVolumeUpload
	ret,err:=helper.HttpPost(httpClient,url,body.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		log.Warn(err)
		return err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return nil
	}else{
		return errors.Error(int32(result))
	}
}

func (this *ProxyHttp) Get(host string,vid int32,key int64,cookie int32) ([]byte,error) {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	jsonSend.Set(vid,"vid")
	jsonSend.Set(key,"key")
	jsonSend.Set(cookie,"cookie")
	url:="http://"+host+constant.CmdVolumeGet
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return nil,errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return nil,err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result!=0 {
		return nil,errors.Error(int32(result))
	}

	img:=jsonRet.Path("img").Data().(string)
	imgData,err:=base64.StdEncoding.DecodeString(img)
	if err!=nil {
		return nil,err
	}else{
		return imgData,nil
	}
}

func (this *ProxyHttp) Delete(host string,vid int32,key int64,replication bool) error {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	jsonSend.Set(vid,"vid")
	jsonSend.Set(key,"key")
	jsonSend.Set(replication,"replication")
	url:="http://"+host+constant.CmdVolumeDelete
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return nil
	}else{
		return errors.Error(int32(result))
	}
}

func (this *ProxyHttp) StorageInit(host string,replication bool) error {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	jsonSend.Set(replication,"replication")
	url:="http://"+host+constant.CmdStoreInit
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return nil
	}else{
		return errors.Error(int32(result))
	}
}

func (this *ProxyHttp) StorageVolumeCompact(host string,vid int32,replication bool) error {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	jsonSend.Set(vid,"vid")
	jsonSend.Set(replication,"replication")
	url:="http://"+host+constant.CmdVolumeCompact
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return nil
	}else{
		return errors.Error(int32(result))
	}
}

func (this *ProxyHttp) StorageVolumeClear(host string,vid int32,replication bool) error {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	jsonSend.Set(vid,"vid")
	jsonSend.Set(replication,"replication")
	url:="http://"+host+constant.CmdVolumeClear
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return nil
	}else{
		return errors.Error(int32(result))
	}
}

func (this *ProxyHttp) StorageGetConfigure(host string) (*gabs.Container,error) {
	jsonSend:=gabs.New()
	jsonSend.Set(helper.UUIDBuild(),"seq")
	url:="http://"+host+constant.CmdStoreConf
	ret,err:=helper.HttpPost(httpClient,url,jsonSend.Bytes())
	if err!=nil {
		log.Warnf("HttpPost[%s] error[%v]",url,err)
		return nil,errors.ErrRpc
	}

	jsonRet,err:=gabs.ParseJSON(ret)
	if err!=nil {
		return nil,err
	}

	result:=jsonRet.Path(helper.CmdResult).Data().(float64)
	if result==0 {
		return jsonRet,nil
	}else{
		return nil,errors.Error(int32(result))
	}
}

func (this *ProxyHttp) Stop() {
	log.Info("proxy http Stop")
}

