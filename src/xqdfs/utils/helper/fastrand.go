package helper

import (
	"sync"
	"time"
)

var(
	rngPool sync.Pool
)

func Uint32() uint32 {
	v := rngPool.Get()
	if v == nil {
		v = &RNG{}
	}
	r := v.(*RNG)
	x := r.Uint32()
	rngPool.Put(r)
	return x
}

func Uint32n(maxN uint32) uint32 {
	x := Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

type RNG struct {
	x uint32
}

func (r *RNG) Uint32() uint32 {
	for r.x == 0 {
		r.x = getRandomUint32()
	}

	x := r.x
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	r.x = x
	return x
}

func (r *RNG) Uint32n(maxN uint32) uint32 {
	x := r.Uint32()
	return uint32((uint64(x) * uint64(maxN)) >> 32)
}

func getRandomUint32() uint32 {
	x := time.Now().UnixNano()
	return uint32((x >> 32) ^ x)
}
