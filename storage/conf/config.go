package conf

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"io/ioutil"
	"encoding/json"

	"xqdfs/storage/needle"

	"github.com/BurntSushi/toml"
	"github.com/Jeffail/gabs"
)

type Config struct {
	NeedleMaxSize int	`json:"needle_max_size"`
	Log		  *Log		`json:"-"`
	Store     *Store	`json:"store"`
	Block     *Block	`json:"block"`
	Index     *Index	`json:"index"`
	Http	  *Http		`json:"http"`
	Dir 	  *Dir		`json:"dir"`
	Configure  *Configure 			`json:"-"`
	Replication	*Replication		`json:"-"`
}

type Log struct {
	Level string
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

type Http struct {
	Host string			`json:"host"`
	Port int			`json:"port"`
}

type Dir struct {
	Path []string		`json:"path"`
	Capacity []int		`json:"capacity"`		//GB
}

type Configure struct {
	Param string
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
func NewConfig(conf string) (c *Config, err error) {
	var (
		file *os.File
		blob []byte
	)
	c = new(Config)
	if file, err = os.Open(conf); err != nil {
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
