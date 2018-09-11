package process

import (
	"io"
	"os"
	"fmt"
	"sync"
	"time"

	"xqdfs/proxy"
	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/helper"
	"xqdfs/storage/store"
	"xqdfs/storage/needle"
	"xqdfs/storage/replication/binlog"
)

type ReplicationTask struct {
	StorageId int32
	StorageAddr string
	Online bool

	s *store.Store
	proxyStorage *proxy.ProxyStorage
	curConfigureFile string	//当前配置文件
	curBinlogFile string 	//当前操作的binlog文件
	curBinlogLine int64		//当前操作的binlog行号

	lock sync.Mutex
	binlogConfigure *binlog.BinlogConfigure
	binlogReader *binlog.BinlogReader
	binlogWriter *binlog.BinlogWriter
	wg sync.WaitGroup
	isRun bool
	signal chan int
}

func NewReplicationTask(path string,s *store.Store,proxyStorage *proxy.ProxyStorage,storageId int32,storageAddr string) (*ReplicationTask,error) {
	log.Debugf("create ReplicationTask [%s]",storageAddr)
	confPath:=path+"/"+helper.Int32ToString(storageId)+".toml"
	binlogPath:=path+"/"+helper.Int32ToString(storageId)+".log"

	binlogConfigure,err:=binlog.NewBinlogConfigure(confPath,storageId,storageAddr)
	if err!=nil {
		log.Error(err)
		return nil,err
	}

	binlogWriter,err:=binlog.NewBinlogWriter(binlogPath)
	if err!=nil {
		log.Error(err)
		return nil,err
	}

	binlogReader,err:=binlog.NewBinlogReader(binlogPath)
	if err!=nil {
		log.Error(err)
		return nil,err
	}

	sync:=&ReplicationTask{
		s:s,
		proxyStorage:proxyStorage,
		StorageId:storageId,
		StorageAddr:storageAddr,
		Online:false,
		curConfigureFile:confPath,
		curBinlogFile:binlogPath,
		curBinlogLine:binlogConfigure.BinlogLine,
		binlogConfigure:binlogConfigure,
		binlogWriter:binlogWriter,
		binlogReader:binlogReader,
		isRun:true,
		signal:make(chan int, 1),
	}
	go sync.start()
	return sync,nil
}

func (this *ReplicationTask) start() {
	for this.isRun {
		func(){
			defer helper.HandleErr()
			err:=this.process()
			if err==io.EOF || err==errors.ErrRpc{
				select {
				case <-time.After(time.Second):
				case <-this.signal:
				}
			}
		}()
	}
	log.Infof("replication task exit [%d][%s]",this.StorageId,this.StorageAddr)
	this.wg.Done()
}

func (this *ReplicationTask) compressFile() {
	if this.curBinlogLine==0 {
		return
	}
	log.Debug("compressFile")
	this.curBinlogLine=0

	//delete
	if this.binlogConfigure !=nil {
		this.binlogConfigure.Close()
	}
	if this.binlogReader !=nil {
		this.binlogReader.Close()
	}
	if this.binlogWriter !=nil {
		this.binlogWriter.Close()
	}

	err := os.Remove(this.curConfigureFile)
	if err!=nil{
		log.Error(err)
	}
	err = os.Remove(this.curBinlogFile)
	if err!=nil{
		log.Error(err)
	}

	//create
	binlogConfigure,err:=binlog.NewBinlogConfigure(this.curConfigureFile,this.StorageId,this.StorageAddr)
	if err!=nil {
		log.Error(err)
	}else{
		this.binlogConfigure=binlogConfigure
	}

	binlogWriter,err:=binlog.NewBinlogWriter(this.curBinlogFile)
	if err!=nil {
		log.Error(err)
	}else{
		this.binlogWriter=binlogWriter
	}

	binlogReader,err:=binlog.NewBinlogReader(this.curBinlogFile)
	if err!=nil {
		log.Error(err)
	}else{
		this.binlogReader=binlogReader
	}
}

func (this *ReplicationTask) Stop() {
	log.Infof("replication task stop [%d]",this.StorageId)

	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.binlogConfigure !=nil {
		this.binlogConfigure.Close()
	}
	if this.binlogReader !=nil {
		this.binlogReader.Close()
	}
	if this.binlogWriter !=nil {
		this.binlogWriter.Close()
	}
}

func (this *ReplicationTask) Destroy() {
	log.Infof("replication task destroy [%d][%s]",this.StorageId,this.StorageAddr)

	this.wg.Add(1)
	this.isRun=false
	this.signal<-1
	this.wg.Wait()
	close(this.signal)

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.binlogConfigure !=nil {
		this.binlogConfigure.Close()
	}
	if this.binlogReader !=nil {
		this.binlogReader.Close()
	}
	if this.binlogWriter !=nil {
		this.binlogWriter.Close()
	}

	err := os.Remove(this.curConfigureFile)
	if err!=nil{
		log.Error(err)
	}
	err = os.Remove(this.curBinlogFile)
	if err!=nil{
		log.Error(err)
	}
}

func (this *ReplicationTask) Write(cmd byte,vid int32,key int64,cookie int32) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.binlogWriter.Write(cmd,vid,key,cookie)
}

func (this *ReplicationTask) process() (err error) {
	var(
		data *binlog.Binlog
	)

	this.lock.Lock()
	data,err=this.binlogReader.Read(this.curBinlogLine+1)
	if err==io.EOF {
		this.compressFile()
	}
	this.lock.Unlock()
	if err!=nil {
		return
	}

	vid:=data.Vid
	key:=data.Key
	cookie:=data.Cookie
	cmd:=data.Cmd
	switch cmd {
	case binlog.BinlogCmdUpload:
		v:= this.s.Volumes[vid]
		if v != nil {
			var n *needle.Needle
			n, err= v.Read(key, cookie)
			if err==nil{
				err=this.proxyStorage.Upload(this.StorageAddr,vid,key,cookie,n.Data,false)
			}
			if err!=nil{
				log.Debug(err)
			}
		}else{
			err=errors.ErrVolumeNotExist
		}
	case binlog.BinlogCmdDelete:
		err=this.proxyStorage.Delete(this.StorageAddr,vid,key,false)
	case binlog.BinlogCmdStorageInit:
		err=this.proxyStorage.StorageInit(this.StorageAddr,false)
	case binlog.BinlogCmdStorageVolumeCompact:
		err=this.proxyStorage.StorageVolumeCompact(this.StorageAddr,vid,false)
	case binlog.BinlogCmdStorageVolumeClear:
		err=this.proxyStorage.StorageVolumeClear(this.StorageAddr,vid,false)
	}

	if err!=errors.ErrRpc{
		this.curBinlogLine++
	}
	return
}

func (this *ReplicationTask) String() string {
	return fmt.Sprintf("%s",this.StorageAddr)
}