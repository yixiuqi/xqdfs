package helper

import (
	"runtime"
	"xqdfs/utils/log"
)

func HandleErr() {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		log.Error(NewStringBuilder().Append("处理中出现异常[").Append(err).Append("]").ToString())
		log.Error(NewStringBuilder().Append("堆栈信息:\n").Append(string(buf[:])).ToString())
	}
}
