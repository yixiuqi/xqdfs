package process

import (
	"xqdfs/storage/replication/binlog"
	"xqdfs/utils/log"
)

type ReplicationDelete struct {
	Vid int32
	Key int64
}

func (this *ReplicationDelete) Process(task map[int32]*ReplicationTask) {
	for _,t:=range task{
		if t.Online {
			err:=t.proxyStorage.Delete(t.StorageAddr,this.Vid,this.Key,false)
			if err==nil{
				return
			}
		}

		err:=t.Write(binlog.BinlogCmdDelete,this.Vid,this.Key,0)
		if err!=nil{
			log.Error(err)
		}
	}
}
