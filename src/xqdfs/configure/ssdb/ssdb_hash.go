package ssdb

import (
	"errors"

	"xqdfs/utils/helper"
	myerr "xqdfs/errors"
)

type SSDBHashItem struct{
	Key string
	Value string
}

type SSDBHash struct{
	connMgr *SSDBConnectMgr
}

func NewSSDBHash(connMgr *SSDBConnectMgr) *SSDBHash {
	return &SSDBHash{
		connMgr:connMgr,
	}
}

func (this *SSDBHash) HSet(name string,key string,value string) error {
	if this.connMgr == nil {
		return errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	resp, err:= link.do("hset", name,key,value)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return errors.New("bad response")
}

func (this *SSDBHash) HGet(name string,key string) (string,error) {
	if this.connMgr == nil {
		return "",errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return "",err
	}
	defer link.close()

	resp, err:= link.do("hget",name,key)
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", myerr.ErrParamNotExist
	}
	return "", errors.New("bad response")
}

func (this *SSDBHash) HSize(name string) (int,error) {
	if this.connMgr == nil {
		return 0,errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return 0,err
	}
	defer link.close()

	resp, err:= link.do("hsize", name)
	if err != nil {
		return 0,err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		count,err:=helper.StringToInt(resp[1])
		if err!=nil{
			return 0,err
		}else{
			return count,nil
		}
	}
	return 0,errors.New("bad response")
}

func (this *SSDBHash) HDel(name string,key string) error {
	if this.connMgr == nil {
		return errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	resp, err:= link.do("hdel", name,key)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return errors.New("bad response")
}

func (this *SSDBHash) HGetAll(name string) ([]SSDBHashItem,error) {
	if this.connMgr == nil {
		return nil,errors.New("connect manager is null")
	}

	link, err := this.connMgr.getConnect()
	if err != nil {
		return nil,err
	}
	defer link.close()

	resp, err:= link.do("hgetall",name)
	if len(resp) >=1 && resp[0] == "ok" {
		count:=len(resp)-1
		if count==0{
			return nil,nil
		}

		items:=make([]SSDBHashItem,count/2)
		j:=1
		for i:=0;i<count/2;i++{
			items[i].Key=resp[j]
			items[i].Value=resp[j+1]
			j+=2
		}
		return items, nil
	}
	if resp[0] == "not_found" {
		return nil, myerr.ErrParamNotExist
	}
	return nil, errors.New("bad response")
}