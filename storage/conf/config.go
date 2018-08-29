package conf

import (
	"time"
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
	"xqdfs/storage/needle"
	"fmt"
)

type Config struct {
	NeedleMaxSize int
	Id string
	Store     *Store
	Block     *Block
	Index     *Index
	Http	  *Http
	Dir 	  *Dir
}

type Store struct {
	VolumeIndex     string
	FreeVolumeIndex string
}

type Block struct {
	BufferSize    int
	SyncWrite     int
	Syncfilerange bool
}

type Index struct {
	BufferSize    int
	SyncWrite     int
	Syncfilerange bool
}

type Http struct {
	Port int
}

type Dir struct {
	Path []string
	Capacity []int	//GB
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// NewConfig new a config.
func NewConfig(conf string) (c *Config, err error) {
	var (
		file *os.File
		blob []byte
	)
	c = new(Config)
	if file, err = os.Open(conf); err != nil {
		return
	}
	if blob, err = ioutil.ReadAll(file); err != nil {
		return
	}
	if err = toml.Unmarshal(blob, c); err == nil {
		c.NeedleMaxSize = needle.Size(c.NeedleMaxSize)
		c.Block.BufferSize = needle.Size(c.Block.BufferSize)
	}
	return
}

func (c *Config) String() string {
	return fmt.Sprintf(`
-----------------------------
NeedleMaxSize[%d]
Store:
VolumeIndex[%s],FreeVolumeIndex[%s]
Block:
BufferSize[%d],SyncWrite[%d],Syncfilerange[%v]
Index:
BufferSize[%d],SyncWrite[%d],Syncfilerange[%v]
Dir:
Path[%v],Capacity[%v]
-----------------------------`,
		c.NeedleMaxSize,
		c.Store.VolumeIndex,c.Store.FreeVolumeIndex,
		c.Block.BufferSize,c.Block.SyncWrite,c.Block.Syncfilerange,
		c.Index.BufferSize,c.Index.SyncWrite,c.Index.Syncfilerange,
		c.Dir.Path,c.Dir.Capacity)
}
