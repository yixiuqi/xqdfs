package conf

import (
	"os"
	"fmt"
	"io/ioutil"

	"xqdfs/utils/conf"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log		  *conf.Log
	Server	  *conf.Server
	Configure  *conf.Configure
}

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
	err = toml.Unmarshal(blob, c)
	return
}

func (c *Config) String() string {
	return fmt.Sprintf(`
-----------------------------
Http:
Port[%d]
-----------------------------`,
		c.Server.Port)
}
