package config

import (
	"testing"
)

func init() {
	//file := "./conf/default.conf"
	//Init(file)
}

func Test_conf(t *testing.T) {
	t.Log(Opt)
	t.Log(Opt.Get("server", "addr"))
}

type Conf struct {
	Log    Section `ego:"log"`
	Server Section `ego:"server"`
}

func TestConfig_Parse(t *testing.T) {
	c := New()
	err := c.Parse("./conf/default.conf")
	if err != nil {
		panic(t)
	}

	t.Log("common mod key:", c.GetKeys("mod"))

	str, _ := c.Get("log_common").String("level")
	t.Log("section:", str)

	conf := Conf{}
	c.Unmarshal(&conf, "ego")
	t.Log(conf)
}
