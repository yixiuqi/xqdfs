package store

import (
	"sync"
	"os"
	"xqdfs/storage/volume"
	"xqdfs/storage/conf"
	"io/ioutil"
	"strings"
	"strconv"
	"fmt"
	"path/filepath"
	myos "xqdfs/storage/os"
	"xqdfs/utils/log"
	"xqdfs/errors"
	"xqdfs/storage/needle"
	"sort"
	"xqdfs/storage/stat"
	"time"
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

// int32Slice deleted offset sort.
type int32Slice []int32
func (p int32Slice) Len() int           { return len(p) }
func (p int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

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
	last_vid	int32
	flock       sync.Mutex // protect FreeId & saveIndex
	vlock       sync.Mutex // protect Volumes map
	Stats 		*stat.Stats
}

func NewStore(c *conf.Config) (s *Store, err error) {
	s = &Store{}
	s.conf = c
	s.FreeId = 0
	s.Stats=&stat.Stats{}
	s.Volumes = make(map[int32]*volume.Volume)
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

func (s *Store) init() (err error) {
	if err = s.parseFreeVolumeIndex(); err == nil {
		err = s.parseVolumeIndex()
	}
	go s.statproc()
	return
}

func (s *Store) statproc() {
	var (
		v    *volume.Volume
		olds *stat.Stats
		news = new(stat.Stats)
	)
	for {
		olds = s.Stats
		*news = *olds
		s.Stats = news // use news instead, for current display
		olds.Reset()
		for _, v = range s.Volumes {
			v.Stats.Calc()
			olds.Merge(v.Stats)
		}
		olds.Calc()
		s.Stats = olds
		time.Sleep(time.Second)
	}
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
		log.Debugf("parse volume index, id: %d, block: %s, index: %s", id, bfile, ifile)
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
	var v *volume.Volume
	if s.vf != nil {
		s.vf.Close()
	}
	if s.fvf != nil {
		s.fvf.Close()
	}
	for _, v = range s.Volumes {
		log.Infof("volume[%d] close", v.Id)
		v.Close()
	}
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
	if ov = s.Volumes[id]; ov != nil {
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

func (s *Store) Write(n *needle.Needle) (vid int32,err error) {
	if len(s.Volumes) == 0 {
		err=errors.ErrStoreNoVolume
		return
	}

	if s.last_vid == 0 {
		id:=make([]int32,0)
		for k,_:= range s.Volumes {
			id=append(id,k)
		}
		sort.Sort(int32Slice(id))
		s.last_vid=id[0]
		log.Debug("set last_vid:",s.last_vid)
	}

	v:= s.Volumes[s.last_vid]
	err = v.Write(n)
	if err == errors.ErrSuperBlockNoSpace {
		id:=make([]int32,0)
		for k,_:= range s.Volumes {
			id=append(id,k)
		}
		sort.Sort(int32Slice(id))
		for _,item:=range id {
			v:= s.Volumes[item]
			err = v.Write(n)
			if err ==nil{
				if s.last_vid!=item{
					s.last_vid=item
					log.Debug("set last_vid:",s.last_vid)
				}
				vid=s.last_vid
				return
			}
		}
	}else{
		vid=s.last_vid
	}
	return
}

func (s *Store) CompactVolume(id int32) (err error) {
	var (
		v, nv      *volume.Volume
		bdir, idir string
	)
	// try check volume
	if v = s.Volumes[id]; v != nil {
		if v.Compact {
			return errors.ErrVolumeInCompact
		}
	} else {
		return errors.ErrVolumeExist
	}
	// find a free volume
	if nv, err = s.freeVolume(id); err != nil {
		return
	}
	log.Infof("start compact volume: (%d) %s to %s", id, v.Block.File, nv.Block.File)
	// no lock here, Compact is no side-effect
	if err = v.StartCompact(nv); err != nil {
		nv.Destroy()
		v.StopCompact(nil)
		return
	}
	s.vlock.Lock()
	if v = s.Volumes[id]; v != nil {
		log.Infof("stop compact volume: (%d) %s to %s", id, v.Block.File, nv.Block.File)
		if err = v.StopCompact(nv); err == nil {
			err = s.saveVolumeIndex()
		}
	} else {
		// never happen
		err = errors.ErrVolumeExist
		log.Errorf("compact volume: %d not exist(may bug)", id)
	}
	s.vlock.Unlock()
	nv.Destroy()
	if err == nil {
		bdir, idir = filepath.Dir(nv.Block.File), filepath.Dir(nv.Indexer.File)
		_, err = s.AddFreeVolume(1, bdir, idir)
	}
	return
}