package binlog

import (
	"os"
	"fmt"
	"io"

	"xqdfs/errors"
	"xqdfs/utils/log"
)

type BinlogWriter struct {
	path string
	f *os.File
}

func NewBinlogWriter(path string) (*BinlogWriter,error) {
	f, err:= os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
		return nil, err
	}

	l,err:=f.Seek(0,io.SeekEnd)
	if err!=nil{
		f.Close()
		log.Errorf("seek(\"%s\") error(%v)", path, err)
		return nil,err
	}else{
		log.Infof("seek end(\"%s\") (%v)", path, l)
	}

	binlog:=&BinlogWriter{
		path:path,
		f:f,
	}
	return binlog,nil
}

func (this *BinlogWriter) Close() {
	if this.f != nil {
		if err:= this.f.Close(); err != nil {
			log.Errorf("file: %s close error(%v)", this.path, err)
		}
		this.f = nil
	}
}

func (this *BinlogWriter) Write(cmd byte,vid int32,key int64,cookie int32) error {
	data:=fmt.Sprintf("%d %05d %019d %010d\n",cmd,vid,key,cookie)
	n,err:=this.f.Write([]byte(data))
	if err!=nil{
		return err
	}

	if n!=BinlogOneLength{
		return errors.ErrBinlogLength
	}

	return nil
}




