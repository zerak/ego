package service

import (
	"sync"

	"github.com/zerak/ego/config"
	"github.com/zerak/ego/net"
)

type DefaultTcpServer struct {
	conf config.Server
}

func (t DefaultTcpServer) Name() string {
	return "DefaultTcpServer"
}

func (t DefaultTcpServer) Init() error {
	return nil
}

func (t DefaultTcpServer) Register(h interface{}) error {
	return nil
}

func (t DefaultTcpServer) Start() error {
	connector := net.NewTcpConnector()
	listener := net.NewTcpListener()
	return listener.ListenAndServe(t.conf.Addr, connector, true)
}

func (t DefaultTcpServer) Stop(group *sync.WaitGroup) {
	group.Done()
}

type DefaultRpcServer struct {
}

func (t DefaultRpcServer) Name() string {
	return "DefaultRpcServer"
}

func (t DefaultRpcServer) Init() error {
	return nil
}

func (t DefaultRpcServer) Register(h interface{}) error {
	return nil
}

func (t DefaultRpcServer) Start() error {
	return nil
}

func (t DefaultRpcServer) Stop(group *sync.WaitGroup) {
	group.Done()
}

// NewTcp new a tcp service
func NewTcp(conf config.Server) *DefaultTcpServer {
	return &DefaultTcpServer{conf: conf}
}

// NewUdp new a udp service
//func NewUdp() *DefaultUdpServer {
//	return &DefaultUdpServer{}
//}

// NewRpc new a rpc service
func NewRpc() *DefaultRpcServer {
	return &DefaultRpcServer{}
}
