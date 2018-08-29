package process

import (
	"xqdfs/storage/replication/binlog"
	"xqdfs/utils/log"
)

type ReplicationStorageVolumeClear struct {
	Vid int32
}

func (this *ReplicationStorageVolumeClear) Process(task map[int32]*ReplicationTask) {
	for _,t:=range task{
		if t.Online {
			err:=t.proxyStorage.StorageVolumeClear(t.StorageAddr,this.Vid,false)
			if err==nil{
				return
			}
		}

		err:=t.Write(binlog.BinlogCmdStorageVolumeClear,this.Vid,0,0)
		if err!=nil{
			log.Error(err)
		}
	}
}
