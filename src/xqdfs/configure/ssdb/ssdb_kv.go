package ssdb

type SSDBKV struct{
}

func NewSSDBKV() *SSDBKV{
	item:=new(SSDBKV)
	return item
}

func (this *SSDBKV) Get(key string) (string,error){
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return "",err
	}
	defer link.close()

	val,err:= link.get(key)
	if err != nil {
		return "",err
	}else{
		return val,nil
	}
}

func (this *SSDBKV) Set(key string,value string) error{
	link, err := SSDBConnectMgrInstance().getConnect()
	if err != nil {
		return err
	}
	defer link.close()

	err= link.set(key,value)
	if err != nil {
		return err
	}else{
		return nil
	}
}
