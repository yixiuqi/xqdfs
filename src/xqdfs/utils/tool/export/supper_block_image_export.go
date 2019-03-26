package main

import (
	"fmt"
	"flag"
	"io/ioutil"

	"xqdfs/utils/log"
	"xqdfs/storage/conf"
	"xqdfs/storage/block"
	"xqdfs/storage/needle"
)

// export -c /home/yimin/xqdfs/bin_storage/xqdfs_storage.toml  -b /home/yimin/xqdfs/data/1_0
func main() {
	var (
		err	error
		configFilePath string
		fileBlockPath string
		config	*conf.Config
		superBlock	*block.SuperBlock
	)

	flag.StringVar(&configFilePath, "c", "./xqdfs_storage.toml", "storage config file path")
	flag.StringVar(&fileBlockPath, "b", "./block", "block file path")
	flag.Parse()

	if config, err = conf.NewConfig(configFilePath); err != nil {
		log.Errorf("NewConfig[%s] error[%v]", configFilePath, err)
		return
	}

	if superBlock,err  = block.NewSuperBlock(fileBlockPath,config);err !=nil {
		log.Errorf("NewSuperBlock[%s] error[%v]", fileBlockPath, err)
		return
	}

	if err = superBlock.Recovery(0, func(rn *needle.Needle, so, eo uint32) (err1 error) {
		if rn.Flag != needle.FlagOK {
			return
		}
		fileName:=fmt.Sprintf("data_image/%v.jpg",rn.Key)
		ioutil.WriteFile(fileName,rn.Data,0644)
		return
	}); err != nil {
		log.Errorf("Recovery() error(%v)", err)
		return
	}
}
