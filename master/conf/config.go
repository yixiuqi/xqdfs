package conf

import (
	"io/ioutil"
	"os"
	"github.com/BurntSushi/toml"
	"fmt"
)

type Config struct {
	Http	  *Http
}

type Http struct {
	Port int
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
	err = toml.Unmarshal(blob, c)
	return
}

func (c *Config) String() string {
	return fmt.Sprintf(`
-----------------------------
Http:
Port[%d]
-----------------------------`,
		c.Http.Port)
}
