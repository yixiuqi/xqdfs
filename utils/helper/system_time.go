package helper

import (
	"time"
)

func CurrentTime() int64 {
	t := time.Now()
	return t.UnixNano() / 1000 / 1000
}

//当前时间转换为字符串
//格式:YYYY-MM-DD hh-mm-ss
func GetTimeString() string {
	tNow := time.Now()
	timeNow := tNow.Format("2006-01-02 15:04:05")
	return timeNow
}
