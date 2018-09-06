package defines

type GroupDalSortById []*GroupDal
func (g GroupDalSortById) Len() int {
	return len(g)
}
func (g GroupDalSortById) Swap(i, j int){
	g[i], g[j] = g[j], g[i]
}
func (g GroupDalSortById) Less(i, j int) bool {
	return g[j].Id > g[i].Id
}
