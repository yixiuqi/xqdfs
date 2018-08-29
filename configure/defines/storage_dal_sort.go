package defines

type StorageDalSortById []*StorageDal
func (g StorageDalSortById) Len() int {
	return len(g)
}
func (g StorageDalSortById) Swap(i, j int){
	g[i], g[j] = g[j], g[i]
}
func (g StorageDalSortById) Less(i, j int) bool {
	return g[j].Id > g[i].Id
}
