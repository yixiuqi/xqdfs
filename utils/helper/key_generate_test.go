package helper

import (
	"testing"
	"fmt"
)

func TestKeyGenerate(t *testing.T) {
	key:=KeyGenerate()
	fmt.Println(key)
	time:=TimeFromKey(key)
	fmt.Println(time)
	fmt.Println(TimeStringFromKey(key))
}
