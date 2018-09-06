package helper

import (
	"os"
	"fmt"
	"time"
	"sync/atomic"
)

var start uint32 = uint32(os.Getpid())
func KeyGenerate() int64 {
	second:= time.Now().UnixNano() / 1000000 / 1000
	key:=atomic.AddUint32(&start,Uint32())
	random:=int64(key) | (second<<32)
	return random
}

func TimeFromKey(key int64) int32 {
	v:=key>>32
	return int32(v)
}

func TimeStringFromKey(key int64) string {
	date_time := time.Unix(int64(TimeFromKey(key)), 0)
	return fmt.Sprintf("%v", date_time)
}
