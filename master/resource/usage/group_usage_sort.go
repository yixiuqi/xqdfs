package usage

type GroupUsageSortById []*GroupUsage
func (g GroupUsageSortById) Len() int {
	return len(g)
}
func (g GroupUsageSortById) Swap(i, j int){
	g[i], g[j] = g[j], g[i]
}
func (g GroupUsageSortById) Less(i, j int) bool {
	return g[j].Id > g[i].Id
}

type GroupUsageSortBySize []*GroupUsage
func (g GroupUsageSortBySize) Len() int {
	return len(g)
}
func (g GroupUsageSortBySize) Swap(i, j int){
	g[i], g[j] = g[j], g[i]
}
func (g GroupUsageSortBySize) Less(i, j int) bool {
	return g[j].Used > g[i].Used
}

type GroupUsageSortByWriteTps []*GroupUsage
func (g GroupUsageSortByWriteTps) Len() int {
	return len(g)
}
func (g GroupUsageSortByWriteTps) Swap(i, j int){
	g[i], g[j] = g[j], g[i]
}
func (g GroupUsageSortByWriteTps) Less(i, j int) bool {
	return g[j].WriteTps > g[i].WriteTps
}
