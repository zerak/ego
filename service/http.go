package service

import (
	"net/http"
	"sync"

	"github.com/zerak/ego/config"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

type DefaultHttpServer struct {
	conf config.Server
}

func (t DefaultHttpServer) Name() string {
	return "DefaultHttpServer"
}

func (t DefaultHttpServer) Init() error {
	return nil
}

func (t DefaultHttpServer) Register(h interface{}) error {
	http.Handle("", h)
	return nil
}

func (t DefaultHttpServer) Start() error {
	return nil
}

func (t DefaultHttpServer) Stop(group *sync.WaitGroup) {
	group.Done()
}

// NewHttp new a http service
func NewHttp() *DefaultHttpServer {
	return &DefaultHttpServer{}
}
