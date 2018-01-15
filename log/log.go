package log

import (
	"fmt"

	"github.com/golang/glog"
)

func Init() {
	glog.CopyStandardLogTo("INFO")
}

func Trace(format string, arg ...interface{}) {
	fmt.Println(format, arg)
}

func Info(format string, arg ...interface{}) {
	//fmt.Printf(format, arg)
	fmt.Println(format, arg)
}

func Debug(format string, arg ...interface{}) {
	fmt.Println(format, arg)
}

func Warn(format string, arg ...interface{}) {
	fmt.Println(format, arg)
}

func Error(format string, arg ...interface{}) {
	fmt.Println(format, arg)
}

func Fatal(format string, arg ...interface{}) {
	fmt.Println(format, arg)
}

func Flush() {
	glog.Flush()
}
