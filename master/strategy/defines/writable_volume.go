package defines

import "fmt"

type WritableVolume struct {
	GroupId int32
	StorageId int32
	StorageAddr string
	VolumeId int32
}

func (this *WritableVolume) String() string {
	return fmt.Sprintf("[%d %d %d]",this.GroupId,this.StorageId,this.VolumeId)
}
