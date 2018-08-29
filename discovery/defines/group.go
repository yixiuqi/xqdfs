package defines

import (
	"encoding/json"

	"xqdfs/utils/stat"
)

type Group struct {
	Id int32				`json:"id"`
	ReadOnly bool			`json:"read_only"`
	Stat *stat.Stats		`json:"stats,omitempty"`
	Storage []*Storage		`json:"storage,omitempty"`
}

func NewGroup() (g *Group) {
	g=new(Group)
	return
}

func (this *Group) String() string {
	j,err:=json.Marshal(this)
	if err!=nil{
		return ""
	}else{
		return string(j)
	}
}

