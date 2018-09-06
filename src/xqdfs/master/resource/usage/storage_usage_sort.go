package usage

//by id 从小到大
type StorageUsageSortById []*StorageUsage
func (s StorageUsageSortById) Len() int {
	return len(s)
}
func (s StorageUsageSortById) Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}
func (s StorageUsageSortById) Less(i, j int) bool {
	return s[j].Id > s[i].Id
}

//by size 从小到大
type StorageUsageSortBySize []*StorageUsage
func (s StorageUsageSortBySize) Len() int {
	return len(s)
}
func (s StorageUsageSortBySize) Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}
func (s StorageUsageSortBySize) Less(i, j int) bool {
	return s[j].Used > s[i].Used
}

//by write tps 从小到大
type StorageUsageSortByWriteTps []*StorageUsage
func (v StorageUsageSortByWriteTps) Len() int {
	return len(v)
}
func (v StorageUsageSortByWriteTps) Swap(i, j int){
	v[i], v[j] = v[j], v[i]
}
func (v StorageUsageSortByWriteTps) Less(i, j int) bool {
	return v[j].WriteTps > v[i].WriteTps
}
