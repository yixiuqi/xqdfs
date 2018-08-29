package binlog

import "fmt"

const(
	BinlogCmdUpload 					= 1
	BinlogCmdDelete 					= 2
	BinlogCmdStorageInit 			= 3
	BinlogCmdStorageVolumeCompact	= 4
	BinlogCmdStorageVolumeClear		= 5

	BinlogOneLength 		= 39
)

type Binlog struct {
	Cmd byte
	Vid int32
	Key int64
	Cookie int32
}

func (this *Binlog) String() string {
	return fmt.Sprintf("cmd[%d] vid[%d] key[%d] cookie[%d]",this.Cmd,this.Vid,this.Key,this.Cookie)
}


