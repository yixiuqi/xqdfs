package defines

import "fmt"

type WritableVolume struct {
	GroupId int32
	StorageId int32
	StorageAddr string
	VolumeId int32
	Free int64
}

func (this *WritableVolume) String() string {
	return fmt.Sprintf("[%d %d %d]",this.GroupId,this.StorageId,this.VolumeId)
}

type WritableVolumeSort []*WritableVolume
func (v WritableVolumeSort) Len() int {
	return len(v)
}
func (v WritableVolumeSort) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v WritableVolumeSort) Less(i, j int) bool {
	return v[j].Free > v[i].Free
}
