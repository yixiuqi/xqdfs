package process

import (
	"xqdfs/utils/log"
	"xqdfs/storage/replication/binlog"
)

type ReplicationStorageVolumeCompact struct {
	Vid int32
}

func (this *ReplicationStorageVolumeCompact) Process(task map[int32]*ReplicationTask) {
	for _,t:=range task{
		if t.Online {
			err:=t.proxyStorage.StorageVolumeClear(t.StorageAddr,this.Vid,false)
			if err==nil{
				return
			}
		}

		err:=t.Write(binlog.BinlogCmdStorageVolumeCompact,this.Vid,0,0)
		if err!=nil{
			log.Error(err)
		}
	}
}
