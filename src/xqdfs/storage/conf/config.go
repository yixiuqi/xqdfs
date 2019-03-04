package conf

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"io/ioutil"
	"encoding/json"

	"xqdfs/utils/conf"
	"xqdfs/storage/needle"

	"github.com/Jeffail/gabs"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Log		  *conf.Log			`json:"-"`
	Server	  *conf.Server		`json:"server"`
	Configure  *conf.Configure	`json:"-"`
	NeedleMaxSize int			`json:"needleMaxSize"`
	Store     *Store			`json:"store"`
	Block     *Block			`json:"block"`
	Index     *Index			`json:"index"`
	Dir 	  *Dir				`json:"dir"`
	Replication	*Replication	`json:"-"`
}

type Store struct {
	VolumeIndex     string	`json:"volume_index"`
	FreeVolumeIndex string	`json:"free_volume_index"`
}

type Block struct {
	BufferSize    int	`json:"buffer_size"`
	SyncWrite     int	`json:"sync_write"`
	Syncfilerange bool	`json:"sync_file_range"`
}

type Index struct {
	BufferSize    int	`json:"buffer_size"`
	SyncWrite     int	`json:"sync_write"`
	Syncfilerange bool	`json:"sync_file_range"`
}

type Dir struct {
	Path []string		`json:"path"`
	Capacity []int		`json:"capacity"`		//GB
}

type Replication struct {
	Path string
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
func NewConfig(configFilePath string) (c *Config, err error) {
	var (
		file *os.File
		blob []byte
	)
	c = new(Config)
	if file, err = os.Open(configFilePath); err != nil {
		return
	}
	defer file.Close()

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

func (c *Config) Json() (*gabs.Container,error) {
	j,err:=json.Marshal(c)
	if err!=nil{
		return nil,err
	}

	dec := json.NewDecoder(bytes.NewBuffer(j))
	dec.UseNumber()
	jsonObj,err:=gabs.ParseJSONDecoder(dec)
	if err!=nil{
		return nil,err
	}else{
		return jsonObj,nil
	}
}
