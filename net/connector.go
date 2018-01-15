package net

import (
	"ego/log"
	"net"
)

type Connector interface {
	// OnConnect a connector call back func
	OnConnect(net.Conn)
}

type TcpConnector struct {
	conn net.Conn
}

func (t *TcpConnector) OnDisconnect() {
	log.Info("%v disconnect", t.conn.RemoteAddr().String())
}

func (t *TcpConnector) OnConnect(conn net.Conn) {
	t.conn = conn
	rs := NewReadStream(conn, NewDefaultPacketHandler())
	rw := NewRWSession(conn, rs, "server id", 100)
	rw.Run(nil, t.OnDisconnect)
}

func NewTcpConnector() *TcpConnector {
	return &TcpConnector{}
}
