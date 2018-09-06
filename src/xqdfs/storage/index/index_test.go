package index

import (
	"xqdfs/storage/conf"
	"testing"
	"os"
	"fmt"
	"xqdfs/utils/log"
)

var (
	testConf = &conf.Config{
		NeedleMaxSize: 4 * 1024 * 1024,
		Index: &conf.Index{
			BufferSize:    4 * 1024 * 1024,
			SyncWrite:     10,
			Syncfilerange: true,
		},
	}
)

func TestIndex(t *testing.T) {
	var (
		i       *Indexer
		err     error
		file    = "./test.idx"
	)
	log.SetLevel("error")
	os.Remove(file)
	defer os.Remove(file)

	{
		fmt.Println("----------------------------------------------------------------------------- test open")
		if i, err = NewIndexer(file, testConf); err != nil {
			t.Errorf("NewIndexer() error(%v)", err)
			t.FailNow()
		}
		i.Close()
		if err = i.Open(); err != nil {
			t.Errorf("Open() error(%v)", err)
			t.FailNow()
		}
		i.Close()
	}

	{
		fmt.Println("----------------------------------------------------------------------------- test write")
		if i, err = NewIndexer(file, testConf); err != nil {
			t.Errorf("NewIndexer() error(%v)", err)
			t.FailNow()
		}
		if err = i.Write(1,1,2); err != nil {
			t.Errorf("Open() error(%v)", err)
			t.FailNow()
		}
		if err = i.Write(2,1,3); err != nil {
			t.Errorf("Open() error(%v)", err)
			t.FailNow()
		}
		i.Close()
	}

	{
		fmt.Println("----------------------------------------------------------------------------- test recovery")
		if i, err = NewIndexer(file, testConf); err != nil {
			t.Errorf("NewIndexer() error(%v)", err)
			t.FailNow()
		}
		defer i.Close()
		if err = i.Recovery(func(ix *Index) error {
			fmt.Println("Key:",ix.Key,"Offset:",ix.Offset,"Size:",ix.Size)
			return nil
		}); err != nil {
			t.Errorf("Recovery() error(%v)", err)
			t.FailNow()
		}
	}
}
