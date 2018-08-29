package conf

import (
	"os"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Log		  *Log		`json:"-"`
	Http	  *Http
	Configure  *Configure
	AllocStrategy *AllocStrategy
}

type Log struct {
	Level string
}

type Http struct {
	Host string
	Port int
}

type Configure struct {
	Param string
}

type AllocStrategy struct {
	//Order
	OrderClearThreshold int		//最少需要多少空闲块
	OrderMinFreeSpace int64		//卷最少需要多少空间
	OrderConsumeCount int 		//选择多少个卷进行随机写
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
		c.Http.Port)
}
