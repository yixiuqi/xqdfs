package process

import (
	"xqdfs/storage/replication/binlog"
	"xqdfs/utils/log"
)

type ReplicationStorageInit struct {

}

func (this *ReplicationStorageInit) Process(task map[int32]*ReplicationTask) {
	for _,t:=range task{
		if t.Online {
			err:=t.proxyStorage.StorageInit(t.StorageAddr,false)
			if err==nil{
				return
			}
		}

		err:=t.Write(binlog.BinlogCmdStorageInit,0,0,0)
		if err!=nil{
			log.Error(err)
		}
	}
}
