package binlog

import (
	"os"
	"fmt"

	"xqdfs/utils/log"
	myos "xqdfs/storage/os"

	"github.com/BurntSushi/toml"
)

type BinlogConfigure struct {
	StorageId int32
	StorageAddr string
	BinlogLine int64	//binlog当前行
	path string 		//配置文件路径
}

func NewBinlogConfigure(path string,id int32,addr string) (*BinlogConfigure,error) {
	binlog:=&BinlogConfigure{
		StorageId:id,
		StorageAddr:addr,
		path:path,
	}

	if myos.Exist(path) == false {
		f, err:= os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
			return nil, err
		}
		defer f.Close()

		enc:=toml.NewEncoder(f)
		err=enc.Encode(binlog)
		if err!=nil{
			return nil,err
		}
	}else{
		f, err:= os.OpenFile(path, os.O_RDONLY, 0664)
		if err != nil {
			log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
			return nil, err
		}
		defer f.Close()

		_,err=toml.DecodeReader(f,binlog)
		if err!=nil{
			return nil,err
		}
	}

	return binlog,nil
}

func (this *BinlogConfigure) String() string {
	return fmt.Sprintf("id[%d] addr[%s] log[%d]",this.StorageId,this.StorageAddr,this.BinlogLine)
}

func (this *BinlogConfigure) Close() {
}

func (this *BinlogConfigure) Update() error {
	f, err:= os.OpenFile(this.path, os.O_WRONLY|os.O_TRUNC, 0664)
	if err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", this.path, err)
		return err
	}
	defer f.Close()

	enc:=toml.NewEncoder(f)
	err=enc.Encode(this)
	if err!=nil{
		return err
	}else{
		return nil
	}
}