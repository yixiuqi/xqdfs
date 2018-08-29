package helper

import (
	"time"
	"os"
	"sync/atomic"
	"fmt"
)

var start uint32 = uint32(os.Getpid())+uint32(time.Now().UnixNano() / 1000000 / 1000)
func KeyGenerate() int64 {
	value:=atomic.AddUint32(&start,1)
	timestamp := time.Now().UnixNano() / 1000000 / 1000
	r:=Uint32() + value
	var v int64
	v=int64(r) | (timestamp<<32)
	return v
}

func TimeFromKey(key int64) int32 {
	v:=key>>32
	return int32(v)
}

func TimeStringFromKey(key int64) string {
	date_time := time.Unix(int64(TimeFromKey(key)), 0)
	return fmt.Sprintf("%v", date_time)
}
