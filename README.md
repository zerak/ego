 EGO
 ===

 [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/zerak/ego/master/LICENSE)
 [![Go Report Card](https://goreportcard.com/badge/github.com/zerak/ego)](https://goreportcard.com/report/github.com/zerak/ego)
 [![Travis branch](https://img.shields.io/travis/zerak/ego/master.svg)](https://travis-ci.org/zerak/ego)
 [![Coverage Status](https://coveralls.io/repos/github/zerak/ego/badge.svg?branch=master)](https://coveralls.io/github/zerak/ego?branch=master)
 [![GoDoc](https://godoc.org/github.com/zerak/ego?status.svg)](https://godoc.org/github.com/zerak/ego)

 `ego` toolset of develop a server


Getting started
---------------

```go
package main

import (
	"github.com/zerak/ego/config"
	"github.com/zerak/ego/log"
	"github.com/zerak/ego/service"
)

func main() {
	log.Info("mahjong logic server info")
	service.Run(service.NewTcp(config.Server{}))
	log.Fatal("mahjong logic server fatal")
}
```

Run it:

```sh
go run main.go
```