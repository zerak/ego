package config

import (
	"fmt"

	"github.com/zerak/goconf"
)

type Server struct {
	conf *goconf.Config

	// devMod dev/pro
	DevMod string

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
	str := fmt.Sprintf("addr:[%v]", s.Addr)
	str += fmt.Sprintf("\nDb:\n")
	for k, v := range s.Db {
		str += fmt.Sprintf("[%v]:[%v]\n", k, v)
	}
	str += fmt.Sprintf("cache:\n")
	for k, v := range s.Cache {
		str += fmt.Sprintf("[%v]:[%v]\n", k, v)
	}
	return str
}

var Opt Server

func Init(path string) {
	Opt.conf = goconf.New()
	err := Opt.conf.Parse(path)
	if err != nil {
		panic(err)
	}
	err = Opt.conf.Unmarshal(&Opt, "ego")
	if err != nil {
		panic(err)
	}
}

func (c *Server) Get(section string, tag string) (string, error) {
	sec := c.conf.Get(section)
	if sec != nil {
		return "", fmt.Errorf("invalid section:%v", section)
	}
	return sec.String(tag)
}

func (c *Server) GetInt(section string, tag string) (int64, error) {
	sec := c.conf.Get(section)
	if sec != nil {
		return 0, fmt.Errorf("invalid section:%v", section)
	}
	return sec.Int(tag)
}
