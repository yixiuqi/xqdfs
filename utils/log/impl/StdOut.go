package impl

import (
	"fmt"
)

type StdOut struct {
	LogOut
}

func NewStdOut() *StdOut {
	item := new(StdOut)
	return item
}

func (this *StdOut) Out(msg string) {
	fmt.Println(msg)
}
