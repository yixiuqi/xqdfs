package defines

import (
	"encoding/json"

	"xqdfs/utils/log"
	"xqdfs/utils/stat"
)

type Storage struct {
	Id int32							`json:"id"`
	Addr string							`json:"addr"`
	Online bool							`json:"online"`
	Stat *stat.Stats					`json:"stats,omitempty"`
	Volume []*Volume					`json:"volumes,omitempty"`
}

func (this *Storage) String() string {
	data,err:=json.Marshal(this)
	if err!=nil{
		log.Warn(err)
		return ""
	}else{
		return string(data)
	}
}
