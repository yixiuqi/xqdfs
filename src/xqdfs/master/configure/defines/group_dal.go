package defines

import (
	"encoding/json"
)

type GroupDal struct {
	Id int32							`json:"id"`
	ReadOnly bool						`json:"read_only"`
	Storage []*StorageDal				`json:"storage,omitempty"`
}

func NewGroupDal() *GroupDal {
	g:=new(GroupDal)
	return g
}

func (this *GroupDal) String() string {
	data,err:=json.Marshal(this)
	if err!=nil {
		return ""
	}else{
		return string(data)
	}
}
