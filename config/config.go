package config

import (
	"fmt"
	"path/filepath"

	"github.com/zerak/goconf"
)

// Server config of server
type Server struct {
	conf *goconf.Config

	// LogRoot the root path of the log
	LogRoot string `ego:"log:root"`

	// LogName the name of the log
	LogName string `ego:"log:name"`

	// LogLevel log level default debug
	LogLevel string `ego:"log:level"`

	// devMod dev/pro
	DevMod string `ego:"server:mod"`

	// [server]
	// addr ip:port
	Addr string `ego:"server:addr"`

	// [db]
	// mysql mysql01=ip:port,mysql02=ip2:port2
	Db map[string]string `ego:"db:mysql:,"`

	// [cache]
	// redis redis01=ip:port,redis02=ip2:port2
	Cache map[string]string `ego:"cache:redis:,"`
}

func (s Server) String() string {
	str := fmt.Sprintf("log path:[%v/%v] level:[%v] ", s.LogRoot, s.LogName, s.LogLevel)
	str += fmt.Sprintf("server addr:[%v] mod:%v ", s.Addr, s.DevMod)
	str += fmt.Sprintf("Db:")
	for k, v := range s.Db {
		str += fmt.Sprintf("[%v]:[%v] ", k, v)
	}
	str += fmt.Sprintf("cache:")
	for k, v := range s.Cache {
		str += fmt.Sprintf("[%v]:[%v] ", k, v)
	}
	return str
}

// Get get section config as string
func (s *Server) Get(section string, tag string) (string, error) {
	sec := s.conf.Get(section)
	if sec == nil {
		return "", fmt.Errorf("invalid section:%v", section)
	}
	return sec.String(tag)
}

// GetInt get section config as int
func (s *Server) GetInt(section string, tag string) (int64, error) {
	sec := s.conf.Get(section)
	if sec == nil {
		return 0, fmt.Errorf("invalid section:%v", section)
	}
	return sec.Int(tag)
}

// Opt the default server config
var Opt Server

// Init load init config file
func Init(path string) {
	Opt.conf = goconf.New()
	absPath, _ := filepath.Abs(path)
	err := Opt.conf.Parse(absPath)
	if err != nil {
		panic(err)
	}
	err = Opt.conf.Unmarshal(&Opt, "ego")
	if err != nil {
		panic(err)
	}
	if Opt.LogRoot == "" {
		Opt.LogRoot = "../"
	}
	if Opt.LogName == "" {
		Opt.LogName = "app"
	}
	if Opt.LogLevel == "" {
		Opt.LogLevel = "debug"
	}
	if Opt.DevMod == "" {
		Opt.DevMod = "dev"
	}
}

// Get get special section config as string
func Get(section string, tag string) (string, error) {
	sec := Opt.conf.Get(section)
	if sec == nil {
		return "", fmt.Errorf("invalid section:%v", section)
	}
	return sec.String(tag)
}

// GetInt get special section config as int
func GetInt(section string, tag string) (int64, error) {
	sec := Opt.conf.Get(section)
	if sec == nil {
		return 0, fmt.Errorf("invalid section:%v", section)
	}
	return sec.Int(tag)
}
