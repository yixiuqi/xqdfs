package impl

import (
	"os"
	"io"
	"time"
)

//写日志到文件
//1.按时间命名文件如:2018_3_7_9_26_3.txt(年_月_日_时_分_秒)
//2.每隔20M生成一个新文件

const(
	LogMaxSize = 1024 * 1024 * 20
)

type FileOut struct{
	channel chan string
	file *os.File
	writeLen int64
}

func NewFileOut() *FileOut{
	f:=&FileOut{
		channel:make(chan string,256),
	}
	go f.process()
	return f
}

func (this *FileOut) Out(msg string){
	this.channel<-msg
}

func (this *FileOut) process(){
	for true {
		item:=<-this.channel
		if item == "" {
			continue
		}

		if this.writeLen > LogMaxSize {
			if this.file!=nil{
				this.file.Close()
				this.file=nil
			}
		}

		if this.file==nil{
			now:= time.Now()
			logFileParentPath  := "log/"+now.Format("20060102");
			logFilePath := logFileParentPath +"/"+now.Format("2006_01_02_15_04_05")+".txt"
			os.MkdirAll(logFileParentPath,0777)

			var f *os.File
			f,err:= os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, 0755)
			if err!=nil{
				continue
			}
			this.file=f
			this.writeLen=0
		}

		l, _ := this.file.Seek(0, io.SeekEnd)
		this.writeLen+=int64(len(item))
		this.file.WriteAt([]byte(item+"\n"), l)
	}

	if this.file!=nil{
		this.file.Close()
	}
}



