package usage

//by id 从小到大
type VolumeUsageSortById []*VolumeUsage
func (v VolumeUsageSortById) Len() int {
	return len(v)
}
func (v VolumeUsageSortById) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v VolumeUsageSortById) Less(i, j int) bool {
	return v[j].Id > v[i].Id
}

//by size 从小到大
type VolumeUsageSortBySize []*VolumeUsage
func (v VolumeUsageSortBySize) Len() int {
	return len(v)
}
func (v VolumeUsageSortBySize) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v VolumeUsageSortBySize) Less(i, j int) bool {
	return v[j].Used > v[i].Used
}

//by write tps 从小到大
type VolumeUsageSortByWriteTps []*VolumeUsage
func (v VolumeUsageSortByWriteTps) Len() int {
	return len(v)
}
func (v VolumeUsageSortByWriteTps) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v VolumeUsageSortByWriteTps) Less(i, j int) bool {
	return v[j].WriteTps > v[i].WriteTps
}

//by util 从小到大
type VolumeUsageSortByUtil []*VolumeUsage
func (v VolumeUsageSortByUtil) Len() int {
	return len(v)
}
func (v VolumeUsageSortByUtil) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v VolumeUsageSortByUtil) Less(i, j int) bool {
	return v[j].Util > v[i].Util
}
