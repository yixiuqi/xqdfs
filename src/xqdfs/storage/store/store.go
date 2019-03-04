package store

import (
	"os"
	"fmt"
	"time"
	"sync"
	"strings"
	"strconv"
	"io/ioutil"
	"path/filepath"

	"xqdfs/errors"
	"xqdfs/utils/log"
	"xqdfs/utils/stat"
	"xqdfs/utils/helper"
	"xqdfs/storage/conf"
	"xqdfs/storage/block"
	"xqdfs/storage/volume"
	myos "xqdfs/storage/os"
)

// Store get all volume meta data from a index file. index contains volume id,
// volume file path, the super block file index ends with ".idx" if the super
// block is /bfs/super_block_1, then the super block index file is
// /bfs/super_block_1.idx.
//
// volume index file format:
//  ---------------------------------
// | block_path,index_path,volume_id |
// | /xxx/block_1,/xxx/block_1.idx,1\n |
// | /xxx/block_2,/xxx/block_2.idx,2\n |
//  ---------------------------------
//
// store -> N volumes
//		 -> volume index -> volume info
//
// volume -> super block -> needle -> photo info
//        -> block index -> needle -> photo info without raw data

const (
	FreeVolumePrefix = "free_block_"
	VolumeIndexExt   = ".idx"
	VolumeFreeId     = -1
)

type Store struct {
	vf          *os.File
	fvf         *os.File
	FreeId      int32
	Volumes     map[int32]*volume.Volume // split volumes lock
	FreeVolumes []*volume.Volume
	conf        *conf.Config
	flock       sync.Mutex // protect FreeId & saveIndex
	vlock       sync.Mutex // protect Volumes map
	Stats 		*stat.Stats
	isRun 		bool
	wg   sync.WaitGroup
	signal      chan int
}

func NewStore(c *conf.Config) (s *Store, err error) {
	s = &Store{
		signal:make(chan int, 1),
		isRun:true,
		conf:c,
		FreeId:0,
		Stats:&stat.Stats{},
		Volumes:make(map[int32]*volume.Volume),
	}
	if s.vf, err = os.OpenFile(c.Store.VolumeIndex, os.O_RDWR|os.O_CREATE|myos.O_NOATIME, 0664); err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", c.Store.VolumeIndex, err)
		s.Close()
		return nil, err
	}
	if s.fvf, err = os.OpenFile(c.Store.FreeVolumeIndex, os.O_RDWR|os.O_CREATE|myos.O_NOATIME, 0664); err != nil {
		log.Errorf("os.OpenFile(\"%s\") error(%v)", c.Store.FreeVolumeIndex, err)
		s.Close()
		return nil, err
	}
	if err = s.init(); err != nil {
		s.Close()
		return nil, err
	}
	return
}

func (s *Store) Init() (err error) {
	path:=s.conf.Dir.Path
	capacity:=s.conf.Dir.Capacity

	if len(path)==0 || len(capacity)==0 {
		log.Error("path or capacity is null")
		err=errors.ErrStoreInitFailed
		return
	}

	pos:=int32(0)
	for p:=0;p<len(path);p++ {
		count:=int(int64(capacity[p])*1024*1024*1024/block.MaxSize)
		log.Infof("storage init path:[%s] count:[%d]",path[p],count)
		for i:=0;i<count;i++{
			pos++

			if v:= s.Volumes[pos]; v != nil {
				err=errors.ErrVolumeExist
				return
			}

			_,err=s.AddFreeVolume(1,path[p],path[p])
			if err!=nil{
				log.Error("storage init error[%v]",err)
				return
			}
			_,err=s.AddVolume(pos)
			if err!=nil{
				log.Error("storage init error[%v]",err)
				return
			}
		}
	}

	for p:=0;p<len(path);p++ {
		_,err=s.AddFreeVolume(1,path[p],path[p])
		if err!=nil{
			log.Error("storage init error[%v]",err)
			return
		}
	}

	return
}

func (s *Store) init() (err error) {
	if err = s.parseFreeVolumeIndex(); err == nil {
		err = s.parseVolumeIndex()
	}
	go s.statsProc()
	go s.volumeStatsUpdateProc()
	return
}

func (s *Store) statsProc() {
	var (
		v    *volume.Volume
		olds *stat.Stats
		news = new(stat.Stats)
	)
	for s.isRun {
		startTime:=helper.CurrentTime()

		olds = s.Stats
		*news = *olds
		s.Stats = news // use news instead, for current display
		olds.Reset()
		s.vlock.Lock()
		for _, v = range s.Volumes {
			v.Stats.Calc()
			olds.Merge(v.Stats)
		}
		s.vlock.Unlock()
		olds.Calc()
		s.Stats = olds

		endTime:=helper.CurrentTime()
		time.Sleep(time.Millisecond*1000-time.Duration(endTime-startTime))
	}
	s.wg.Done()
	log.Debug("Store.statsProc exit")
}

func (s *Store) volumeStatsUpdateProc() {
	for s.isRun {
		volumes:=make(map[int32]*volume.Volume)
		s.vlock.Lock()
		for k,v:= range s.Volumes {
			volumes[k]=v
		}
		s.vlock.Unlock()

		for _,v:= range volumes {
			v.StoreStats()
		}
		log.Debug("store stats")
		select {
		case <-time.After(time.Second * 60):
		case <-s.signal:
		}
	}
	s.wg.Done()
	log.Debug("Store.volumeStatsUpdateProc exit")
}

func (s *Store) parseIndex(lines []string) (im map[int32]struct{}, ids []int32, bfs, ifs []string, err error) {
	var (
		id    int64
		vid   int32
		line  string
		bfile string
		ifile string
		seps  []string
	)
	im = make(map[int32]struct{})
	for _, line = range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if seps = strings.Split(line, ","); len(seps) != 3 {
			log.Errorf("volume index: \"%s\" format error", line)
			err = errors.ErrStoreVolumeIndex
			return
		}
		bfile = seps[0]
		ifile = seps[1]
		if id, err = strconv.ParseInt(seps[2], 10, 32); err != nil {
			log.Errorf("volume index: \"%s\" format error[%v]", line,err)
			return
		}
		vid = int32(id)
		ids = append(ids, vid)
		bfs = append(bfs, bfile)
		ifs = append(ifs, ifile)
		im[vid] = struct{}{}
	}
	return
}

func (s *Store) parseVolumeIndex() (err error) {
	var (
		i          int
		ok         bool
		id         int32
		bfile      string
		ifile      string
		v          *volume.Volume
		data       []byte
		ids       []int32
		lines      []string
		bfs, ifs []string
	)
	if data, err = ioutil.ReadAll(s.vf); err != nil {
		log.Errorf("ioutil.ReadAll() error(%v)", err)
		return
	}
	lines = strings.Split(string(data), "\n")
	if _, ids, bfs, ifs, err = s.parseIndex(lines); err != nil {
		return
	}

	for i = 0; i < len(bfs); i++ {
		id, bfile, ifile = ids[i], bfs[i], ifs[i]
		if _, ok = s.Volumes[id]; ok {
			continue
		}
		if v, err = newVolume(id, bfile, ifile, s.conf); err != nil {
			return
		}
		s.Volumes[id] = v
	}
	return
}

func (s *Store) parseFreeVolumeIndex() (err error) {
	var (
		i     int
		id    int32
		bfile string
		ifile string
		v     *volume.Volume
		data  []byte
		ids   []int32
		lines []string
		bfs   []string
		ifs   []string
	)
	if data, err = ioutil.ReadAll(s.fvf); err != nil {
		log.Errorf("ioutil.ReadAll() error(%v)", err)
		return
	}
	lines = strings.Split(string(data), "\n")
	if _, ids, bfs, ifs, err = s.parseIndex(lines); err != nil {
		return
	}
	for i = 0; i < len(bfs); i++ {
		id, bfile, ifile = ids[i], bfs[i], ifs[i]
		if v, err = newVolume(id, bfile, ifile, s.conf); err != nil {
			return
		}
		v.Close()
		s.FreeVolumes = append(s.FreeVolumes, v)
		if id = s.fileFreeId(bfile); id > s.FreeId {
			s.FreeId = id
		}
	}
	log.Info("FreeId:",s.FreeId)
	return
}

func (s *Store) file(id int32, bdir, idir string, i int) (bfile, ifile string) {
	var file = fmt.Sprintf("%d_%d", id, i)
	bfile = filepath.Join(bdir, file)
	file = fmt.Sprintf("%d_%d%s", id, i, VolumeIndexExt)
	ifile = filepath.Join(idir, file)
	return
}

func (s *Store) freeFile(id int32, bdir, idir string) (bfile, ifile string) {
	var file = fmt.Sprintf("%s%d", FreeVolumePrefix, id)
	bfile = filepath.Join(bdir, file)
	file = fmt.Sprintf("%s%d%s", FreeVolumePrefix, id, VolumeIndexExt)
	ifile = filepath.Join(idir, file)
	return
}

func (s *Store) fileFreeId(file string) (id int32) {
	var (
		fid    int64
		fidStr = strings.Replace(filepath.Base(file), FreeVolumePrefix, "", -1)
	)
	fid, _ = strconv.ParseInt(fidStr, 10, 32)
	id = int32(fid)
	return
}

func (s *Store) Close() {
	log.Info("store close")
	s.wg.Add(2)
	s.isRun=false
	s.signal<-1
	s.wg.Wait()

	var v *volume.Volume
	if s.vf != nil {
		s.vf.Close()
	}
	if s.fvf != nil {
		s.fvf.Close()
	}
	s.vlock.Lock()
	for _, v = range s.Volumes {
		log.Infof("volume[%d] close", v.Id)
		v.Close()
	}
	s.vlock.Unlock()
	return
}

func (s *Store) saveVolumeIndex() (err error) {
	var (
		tn, n int
		v     *volume.Volume
	)
	if _, err = s.vf.Seek(0, os.SEEK_SET); err != nil {
		log.Errorf("vf.Seek() error(%v)", err)
		return
	}
	for _, v = range s.Volumes {
		if n, err = s.vf.WriteString(fmt.Sprintf("%s\n", string(v.Meta()))); err != nil {
			log.Errorf("vf.WriteString() error(%v)", err)
			return
		}
		tn += n
	}
	if err = s.vf.Sync(); err != nil {
		log.Errorf("vf.Sync() error(%v)", err)
		return
	}
	if err = os.Truncate(s.conf.Store.VolumeIndex, int64(tn)); err != nil {
		log.Errorf("os.Truncate() error(%v)", err)
	}
	return
}

func (s *Store) saveFreeVolumeIndex() (err error) {
	var (
		tn, n int
		v     *volume.Volume
	)
	if _, err = s.fvf.Seek(0, os.SEEK_SET); err != nil {
		log.Errorf("fvf.Seek() error(%v)", err)
		return
	}
	for _, v = range s.FreeVolumes {
		if n, err = s.fvf.WriteString(fmt.Sprintf("%s\n", string(v.Meta()))); err != nil {
			log.Errorf("fvf.WriteString() error(%v)", err)
			return
		}
		tn += n
	}
	if err = s.fvf.Sync(); err != nil {
		log.Errorf("fvf.saveFreeVolumeIndex Sync() error(%v)", err)
		return
	}
	if err = os.Truncate(s.conf.Store.FreeVolumeIndex, int64(tn)); err != nil {
		log.Errorf("os.Truncate() error(%v)", err)
	}
	return
}

func newVolume(id int32, bfile, ifile string, c *conf.Config) (v *volume.Volume, err error) {
	if v, err = volume.NewVolume(id, bfile, ifile, c); err != nil {
		log.Errorf("newVolume(%d, %s, %s) error(%v)", id, bfile, ifile, err)
	}
	return
}

func (s *Store) addVolume(id int32, nv *volume.Volume) {
	var (
		vid     int32
		v       *volume.Volume
		volumes = make(map[int32]*volume.Volume, len(s.Volumes)+1)
	)
	for vid, v = range s.Volumes {
		volumes[vid] = v
	}
	volumes[id] = nv
	// goroutine safe replace
	s.Volumes = volumes
}

func (s *Store) AddVolume(id int32) (v *volume.Volume, err error) {
	var ov *volume.Volume
	// try check exists
	s.vlock.Lock()
	ov = s.Volumes[id]
	s.vlock.Unlock()
	if ov != nil {
		return nil, errors.ErrVolumeExist
	}
	// find a free volume
	if v, err = s.freeVolume(id); err != nil {
		return
	}
	s.vlock.Lock()
	if ov = s.Volumes[id]; ov == nil {
		s.addVolume(id, v)
		err = s.saveVolumeIndex()
		if err != nil {
			log.Errorf("add volume: %d error(%v), local index or zookeeper index may save failed", id, err)
		}
	} else {
		err = errors.ErrVolumeExist
	}
	s.vlock.Unlock()
	if err == errors.ErrVolumeExist{
		v.Destroy()
	}
	return
}

func (s *Store) AddFreeVolume(n int, bdir, idir string) (sn int, err error) {
	var (
		i            int
		bfile, ifile string
		v            *volume.Volume
	)
	s.flock.Lock()
	for i = 0; i < n; i++ {
		s.FreeId++
		bfile, ifile = s.freeFile(s.FreeId, bdir, idir)
		if myos.Exist(bfile) || myos.Exist(ifile) {
			continue
		}
		if v, err = newVolume(VolumeFreeId, bfile, ifile, s.conf); err != nil {
			// if no free space, delete the file
			os.Remove(bfile)
			os.Remove(ifile)
			break
		}
		v.Close()
		s.FreeVolumes = append(s.FreeVolumes, v)
		sn++
		log.Infof("create free volume bfile[%s] ifile[%s]",bfile,ifile)
	}
	err = s.saveFreeVolumeIndex()
	s.flock.Unlock()
	return
}

func (s *Store) freeVolume(id int32) (v *volume.Volume, err error) {
	var (
		i                                        int
		bfile, nbfile, ifile, nifile, bdir, idir string
	)
	s.flock.Lock()
	defer s.flock.Unlock()
	if len(s.FreeVolumes) == 0 {
		err = errors.ErrStoreNoFreeVolume
		return
	}
	v = s.FreeVolumes[0]
	s.FreeVolumes = s.FreeVolumes[1:]
	v.Id = id
	bfile, ifile = v.Block.File, v.Indexer.File
	bdir, idir = filepath.Dir(bfile), filepath.Dir(ifile)
	for {
		nbfile, nifile = s.file(id, bdir, idir, i)
		if !myos.Exist(nbfile) && !myos.Exist(nifile) {
			break
		}
		i++
	}
	log.Infof("rename block: %s to %s", bfile, nbfile)
	log.Infof("rename index: %s to %s", ifile, nifile)
	if err = os.Rename(ifile, nifile); err != nil {
		log.Errorf("os.Rename(\"%s\", \"%s\") error(%v)", ifile, nifile, err)
		v.Destroy()
		return
	}
	if err = os.Rename(bfile, nbfile); err != nil {
		log.Errorf("os.Rename(\"%s\", \"%s\") error(%v)", bfile, nbfile, err)
		v.Destroy()
		return
	}
	v.Block.File = nbfile
	v.Indexer.File = nifile
	if err = v.Open(); err != nil {
		v.Destroy()
		return
	}
	err = s.saveFreeVolumeIndex()
	return
}

func (s *Store) CompactVolume(id int32) (err error) {
	var (
		v, nv      *volume.Volume
		bdir, idir string
	)
	// try check volume
	s.vlock.Lock()
	v = s.Volumes[id]
	s.vlock.Unlock()
	if v != nil {
		if v.Compact {
			return errors.ErrVolumeInCompact
		}
	} else {
		return errors.ErrVolumeNotExist
	}
	// find a free volume
	if nv, err = s.freeVolume(id); err != nil {
		return
	}
	log.Infof("start compact volume: (%d) %s to %s", id, v.Block.File, nv.Block.File)
	// no lock here, Compact is no side-effect
	if err = v.StartCompact(nv); err != nil {
		log.Error(err)
		nv.Destroy()
		v.StopCompact(nil)
		return
	}
	s.vlock.Lock()
	if v = s.Volumes[id]; v != nil {
		log.Infof("stop compact volume: (%d) %s to %s", id, v.Block.File, nv.Block.File)
		if err = v.StopCompact(nv); err == nil {
			err = s.saveVolumeIndex()
		}else{
			log.Warn(err)
		}
	} else {
		// never happen
		err = errors.ErrVolumeNotExist
		log.Errorf("compact volume: %d not exist(may bug)", id)
	}
	s.vlock.Unlock()
	nv.Destroy()
	if err == nil {
		bdir, idir = filepath.Dir(nv.Block.File), filepath.Dir(nv.Indexer.File)
		_, err = s.AddFreeVolume(1, bdir, idir)
	}else{
		log.Warn(err)
	}
	return
}