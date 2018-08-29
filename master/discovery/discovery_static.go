package discovery

type Discovery_Static struct {
	kv []*KVPair
}

func NewDiscovery_Static(kv []*KVPair) *Discovery_Static{
	item:=new(Discovery_Static)
	item.kv=kv
	return item
}

func (this *Discovery_Static) GetServices() []*KVPair {
	return this.kv
}

func (this *Discovery_Static) WatchService() chan []*KVPair {
	return nil
}

func (this *Discovery_Static) RemoveWatcher(ch chan []*KVPair){
}

func (this *Discovery_Static) IsAvailable() bool{
	return true
}

func (this *Discovery_Static) Destroy(){
}
