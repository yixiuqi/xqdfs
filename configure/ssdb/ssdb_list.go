package ssdb

import (
	"fmt"

	"xqdfs/utils/helper"
)

type SSDBList struct{
}

func NewSSDBList() *SSDBList{
	item:=new(SSDBList)
	return item
}

func (this *SSDBList) QPush(name string,item string) error{
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	resp, err:= link.do("qpush", name,item)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return fmt.Errorf("bad response")
}

func (this *SSDBList) QPop(name string) (string,error){
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return "",err
	}
	defer link.close()

	resp, err:= link.do("qpop", name)
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", fmt.Errorf("not_found")
	}
	return "", fmt.Errorf("bad response")
}

func (this *SSDBList) QFront(name string) (string,error){
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return "",err
	}
	defer link.close()

	resp, err:= link.do("qfront", name)
	if len(resp) == 2 && resp[0] == "ok" {
		return resp[1], nil
	}
	if resp[0] == "not_found" {
		return "", fmt.Errorf("not_found")
	}
	return "", fmt.Errorf("bad response")
}

func (this *SSDBList) QSize(name string) (int64,error){
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return 0,err
	}
	defer link.close()

	resp, err:= link.do("qsize", name)
	if len(resp) == 2 && resp[0] == "ok" {
		size,_:=helper.StringToInt64(resp[1])
		return size, nil
	}
	if resp[0] == "not_found" {
		return 0, fmt.Errorf("not_found")
	}
	return 0, fmt.Errorf("bad response")
}

func (this *SSDBList) QClear(name string) error{
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	resp, err:= link.do("qclear", name)
	if err != nil {
		return err
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return fmt.Errorf("bad response")
}
