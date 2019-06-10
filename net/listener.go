package net

import (
	"net"
)

type Listener interface {
	// ListenAndServe listen on addrStr and serve
	ListenAndServe(addrStr string, handler Connector) error
}

type TcpListener struct {
}

func (t *TcpListener) listen(addrStr string) (net.Listener, error) {
	addr, err := net.ResolveTCPAddr("tcp4", addrStr)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", addr)
}

func (t *TcpListener) serve(listener net.Listener, handler Connector, async bool) error {
	serveFunc := func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handler.OnConnect(conn)
		}
	}
	if async {
		go serveFunc()
	} else {
		serveFunc()
	}
	return nil
}

func (t *TcpListener) ListenAndServe(addrStr string, handler Connector, async bool) error {
	listener, err := t.listen(addrStr)
	if err == nil {
		err = t.serve(listener, handler, async)
	}
	return err
}

func NewTcpListener() *TcpListener {
	return &TcpListener{}
}

type UdpListener struct {
}

func (u *UdpListener) listen(addrStr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	return conn, err
}

func (u *UdpListener) ListenAndServe(addrStr string, handler func(*net.UDPConn), async bool) error {
	conn, err := u.listen(addrStr)
	if err != nil {
		return err
	}

	serveFunc := func() {
		defer conn.Close()
		for {
			handler(conn)
		}
	}
	if async {
		go serveFunc()
	} else {
		serveFunc()
	}
	return nil
}

func NewUdpListener() *UdpListener {
	return &UdpListener{}
}
