package process

import (
	"xqdfs/utils/log"
	"xqdfs/storage/replication/binlog"
)

type ReplicationUpload struct {
	Vid int32
	Key int64
	Cookie int32
	Image []byte
}

func (this *ReplicationUpload) Process(task map[int32]*ReplicationTask) {
	for _,t:=range task{
		if t.Online {
			err:=t.proxyStorage.Upload(t.StorageAddr,this.Vid,this.Key,this.Cookie,this.Image,false)
			if err==nil{
				return
			}
		}

		err:=t.Write(binlog.BinlogCmdUpload,this.Vid,this.Key,this.Cookie)
		if err!=nil{
			log.Error(err)
		}
	}
}
