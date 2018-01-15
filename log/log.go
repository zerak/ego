package log

import (
	"github.com/zerak/log"

	"github.com/zerak/ego/config"
)

func Init() {
	level, _ := log.ParseLevel(config.Opt.LogLevel)
	defer log.Uninit(log.InitMultiFileAndColoredConsole(config.Opt.LogRoot, config.Opt.LogName, level))
	log.SetLevel(level)
}

func Trace(format string, arg ...interface{}) {
	log.Trace(format, arg...)
}

func Info(format string, arg ...interface{}) {
	log.Info(format, arg...)
}

func Debug(format string, arg ...interface{}) {
	log.Debug(format, arg...)
}

func Warn(format string, arg ...interface{}) {
	log.Warn(format, arg...)
}

func Error(format string, arg ...interface{}) {
	log.Error(format, arg...)
}

func Fatal(format string, arg ...interface{}) {
	log.Fatal(format, arg...)
}
