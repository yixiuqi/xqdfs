package binlog

import (
	"os"
	"strings"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
)

type BinlogReader struct {
	path string
	f *os.File
}

func NewBinlogReader(path string) (*BinlogReader,error) {
	f, err:= os.OpenFile(path, os.O_RDONLY, 0664)
	if err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", path, err)
		return nil, err
	}

	binlog:=&BinlogReader{
		path:path,
		f:f,
	}
	return binlog,nil
}

func (this *BinlogReader) Close() {
	if this.f != nil {
		if err:= this.f.Close(); err != nil {
			log.Errorf("file: %s close error(%v)", this.path, err)
		}
		this.f = nil
	}
}

func (this *BinlogReader) Read(line int64) (*Binlog,error) {
	offset:=(line-1)*BinlogOneLength
	data:=make([]byte,BinlogOneLength)
	n,err:=this.f.ReadAt(data,offset)
	if err!=nil {
		return nil,err
	}
	if n!=len(data){
		log.Error("n!=BinlogOneLength")
		return nil,errors.ErrBinlogRead
	}
	if data[BinlogOneLength-1]!='\n'{
		log.Error("last is not \n")
		return nil,errors.ErrBinlogRead
	}
	data=data[0:BinlogOneLength-1]

	str:=string(data)
	items:=strings.Split(str," ")
	if len(items)!=4 {
		log.Error("item count is not 4")
		return nil,errors.ErrBinlogRead
	}

	cmd,err:=helper.StringToInt32(items[0])
	if err!=nil{
		return nil,err
	}
	vid,err:=helper.StringToInt32(items[1])
	if err!=nil{
		return nil,err
	}
	key,err:=helper.StringToInt64(items[2])
	if err!=nil{
		return nil,err
	}
	cookie,err:=helper.StringToInt32(items[3])
	if err!=nil{
		return nil,err
	}

	binlog:=&Binlog{
		Cmd:byte(cmd),
		Vid:int32(vid),
		Key:key,
		Cookie:cookie,
	}
	return binlog,nil
}
