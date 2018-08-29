package discovery

type KVPair struct {
	Key   string
	Value string
}

type Discovery interface {
	GetServices() []*KVPair
	WatchService() chan []*KVPair
	RemoveWatcher(ch chan []*KVPair)
	IsAvailable() bool
	Destroy()
}