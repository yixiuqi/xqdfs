package helper

import (
	"strings"
	"net"
	"regexp"
)

func HostAddrCheck(addr string) bool {
	items:=strings.Split(addr,":")
	if items==nil||len(items)!=2{
		return false
	}

	a:=net.ParseIP(items[0])
	if a==nil{
		return false
	}

	match, err := regexp.MatchString("^[0-9]*$", items[1])
	if err!=nil{
		return false
	}

	if match==false{
		return false
	}

	return true
}
