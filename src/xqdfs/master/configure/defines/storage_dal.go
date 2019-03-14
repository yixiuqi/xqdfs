package defines

import (
	"encoding/json"
)

type StorageDal struct {
	Id int32		`json:"id"`
	Addr string		`json:"addr"`
	Desc string		`json:"desc,omitempty"`
}

func NewStorageDal() *StorageDal {
	s:=new(StorageDal)
	return s
}

func (this *StorageDal) String() string {
	data,err:=json.Marshal(this)
	if err!=nil {
		return ""
	}else{
		return string(data)
	}
}