package config

import (
	"testing"
)

func init() {
	file := "./conf/default.conf"
	Init(file)
}

func Test_conf(t *testing.T) {
	t.Log(Opt)
	t.Log(Opt.Get("server", "addr"))
}
