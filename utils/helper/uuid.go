package helper

import (
	"strings"

	"github.com/satori/go.uuid"
)

func UUIDBuild() string {
	u, _ := uuid.NewV4()
	str:=strings.Replace(u.String(),"-","",-1)
	return str
}
