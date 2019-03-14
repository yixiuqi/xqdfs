package ssdb

import "errors"

type SSDBKV struct{
	connMgr *SSDBConnectMgr
}

func NewSSDBKV(connMgr *SSDBConnectMgr) *SSDBKV {
	return &SSDBKV{
		connMgr:connMgr,
	}
}

func (this *SSDBKV) Get(key string) (string,error){
	if this.connMgr == nil {
		return "",errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return "",err
	}
	defer link.close()

	val,err:= link.get(key)
	if err != nil {
		return "",err
	}else{
		return val,nil
	}
}

func (this *SSDBKV) Set(key string,value string) error {
	if this.connMgr == nil {
		return errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	err= link.set(key,value)
	if err != nil {
		return err
	}else{
		return nil
	}
}

func (this *SSDBKV) Setx(key string,value string,ttl int) error {
	if this.connMgr == nil {
		return errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	err= link.setx(key,value,ttl)
	if err != nil {
		return err
	}else{
		return nil
	}
}
