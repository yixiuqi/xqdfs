package helper

import (
	"os"
	"time"
	"math/rand"
	"sync/atomic"
)

var start int64 = int64(os.Getpid())
func KeyGenerate() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random:=r.Int63()+atomic.AddInt64(&start,1)
	if random<0{
		random=random*-1
	}
	return random
}
