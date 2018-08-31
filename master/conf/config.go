package conf

import (
	"os"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log		  *Log
	Server	  *Server
	Configure  *Configure
}

type Log struct {
	Level string
}

type Server struct {
	Host string
	Port int
}

type Configure struct {
	Param string
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
