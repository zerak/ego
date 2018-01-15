package service

import (
	"flag"
	"os"
	"sync"

	"ego/config"
	"ego/log"
	"ego/signal"
)

type Servicer interface {
	// Name the name of the service
	Name() string

	// Init init service
	Init() error

	// Register register service handler
	Register(h interface{}) error

	// Start start a service no blocking
	Start() error

	// Stop wait for all job done
	// then call sync.WaitGroup.Done
	Stop(*sync.WaitGroup)
}

func Run(services ...Servicer) {
	// 1 init config module
	confFile := flag.String("c", "./conf/default.conf", " default config file path")
	flag.Parse()
	config.Init(*confFile)

	// 2 init log module
	log.Init()

	log.Info("run services")
	for _, s := range services {
		err := s.Init()
		if err != nil {
			log.Fatal("service:%v init err:%v", s.Name(), err)
		}
		log.Info("service:%v init ok", s.Name())
	}

	for _, s := range services {
		err := s.Start()
		if err != nil {
			log.Fatal("service:%v start err:%v", s.Name(), err)
		}
		log.Info("service:%v start ok:%v", s.Name(), err)
	}

	// wait signal
	// Ctrl+C or kill -p
	signal.Wait(os.Interrupt)

	var wg sync.WaitGroup
	for _, s := range services {
		wg.Add(1)
		s.Stop(&wg)
		log.Info("service:%v stop", s.Name())
	}
	wg.Wait()
	log.Info("all services exit")
}
