package defines

type Discovery interface {
	Groups() []*Group
	Storages() []*Storage
	Stop()
}
