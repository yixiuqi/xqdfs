package log

import (
	"bytes"
	"container/list"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"xqdfs/utils/log/impl"
)

/*
	all- 开启所有日志
	trace- 更详细的跟踪信息
	debug- 调试信息，可记录详细的业务处理到哪一步了，以及当前的变量状态。
	info- 有意义的事件信息,如程序启动、关闭事件、收到请求事件等。
	warn- 警告信息，如程序调用了一个即将作废的接口、接口的不当使用、运行状态不是期望的但仍可继续处理等。
	error- 严重的、造成服务中断的错误。
	off- 关闭所有日志
*/

var level int32 = LevelAll
var items *list.List

func init() {
	items = list.New()
	items.PushBack(impl.NewStdOut())
}

func systemInfo() string {
	funcName, file, line, ok := runtime.Caller(2)
	str := bytes.Buffer{}
	if ok {
		funcName = funcName
		//str.WriteString(runtime.FuncForPC(funcName).Name())
		//str.WriteString(" ")
		if len(file)<48{
			str.WriteString(file)
		}else{
			file = string([]byte(file)[len(file)-48:len(file)])
			str.WriteString(".."+file)
		}
		str.WriteString(" ")
		str.WriteString(strconv.Itoa(line))
	}
	return str.String()
}

func buildMessage(level string, msg string, call string) string {
	time := time.Now().Format("2006-01-02 15:04:05")
	str := bytes.Buffer{}
	str.WriteString("[")
	str.WriteString(level)
	str.WriteString("] ")
	str.WriteString(time)
	str.WriteString(" [")
	str.WriteString(call)
	str.WriteString("]")
	str.WriteString(" ")
	str.WriteString(msg)
	return str.String()
}

func out(msg string) {
	for e := items.Front(); e != nil; e = e.Next() {
		e.Value.(impl.LogOut).Out(msg)
	}
}

func AppendOut(out impl.LogOut) {
	items.PushBack(out)
}

func SetLevel(l string) {
	switch l {
	case "all":
		level = LevelAll
	case "off":
		level = LevelOff
	case "trace":
		level = LevelTrace
	case "debug":
		level = LevelDebug
	case "warn":
		level = LevelWarn
	case "info":
		level = LevelInfo
	case "error":
		level = LevelError
	default:
		level = LevelDebug
	}
}

func Trace(a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprint(a...)
	if level >= LevelTrace {
		out(buildMessage("trace", msg, str))
	}
}

func Tracef(format string, a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprintf(format,a...)
	if level >= LevelTrace {
		out(buildMessage("trace", msg, str))
	}
}

func Debug(a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprint(a...)
	if level >= LevelDebug {
		out(buildMessage("debug", msg, str))
	}
}

func Debugf(format string, a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprintf(format, a...)
	if level >= LevelDebug {
		out(buildMessage("debug", msg, str))
	}
}

func Info(a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprint(a...)
	if level >= LevelInfo {
		out(buildMessage("info", msg, str))
	}
}

func Infof(format string, a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprintf(format,a...)
	if level >= LevelInfo {
		out(buildMessage("info", msg, str))
	}
}

func Warn(a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprint(a...)
	if level >= LevelWarn {
		out(buildMessage("warn", msg, str))
	}
}

func Warnf(format string, a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprintf(format,a...)
	if level >= LevelWarn {
		out(buildMessage("warn", msg, str))
	}
}

func Error(a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprint(a...)
	if level >= LevelError {
		out(buildMessage("error", msg, str))
	}
}

func Errorf(format string, a ...interface{}) {
	str := systemInfo()
	msg := fmt.Sprintf(format,a...)
	if level >= LevelError {
		out(buildMessage("error", msg, str))
	}
}
