package ssdb

import "testing"

func TestConfigureSSDB_ParamSetx(t *testing.T) {
	ssdb:=NewConfigureSSDB("192.168.10.90:18888")
	ssdb.ParamSetx("a","a1",10)
}
