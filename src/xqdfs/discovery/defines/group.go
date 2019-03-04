package defines

import (
	"encoding/json"

	"xqdfs/utils/log"
	"xqdfs/utils/stat"
)

type Group struct {
	Id int32				`json:"id"`
	ReadOnly bool			`json:"read_only"`
	Stat *stat.Stats		`json:"stats,omitempty"`
	Storage []*Storage		`json:"storage,omitempty"`
}

func (this *Group) String() string {
	j,err:=json.Marshal(this)
	if err!=nil{
		log.Warn(err)
		return ""
	}else{
		return string(j)
	}
}

