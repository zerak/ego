package net

import (
	"net"
	"sync/atomic"
	"time"

	buf "github.com/zerak/ego/buffer"
	"github.com/zerak/ego/log"
)

type EncryptFunc func(dst, src []byte)
type DecryptFunc func(dst, src []byte)

type Session interface {
	Id() string
	Send(*buf.Buffer)
	Run(onNewSession, onQuitSession func())
	Quit()
	SetEncrypt(encrypt EncryptFunc)
	SetDecrypt(decrypt DecryptFunc)
}

// Write Session
type WSession struct {
	conn   net.Conn
	id     string
	closed int32

	writeQuit chan struct{}
	writeChan chan *buf.Buffer

	encryptBuf []byte

	encrypt     EncryptFunc
	encryptChan chan EncryptFunc
}

func (ws *WSession) Id() string                     { return ws.id }
func (ws *WSession) setClosed()                     { atomic.StoreInt32(&ws.closed, 1) }
func (ws *WSession) getClosed() bool                { return atomic.LoadInt32(&ws.closed) == 1 }
func (ws *WSession) SetEncrypt(encrypt EncryptFunc) { ws.encryptChan <- encrypt }

func (ws *WSession) Send(b *buf.Buffer) {
	if b.Len() > 0 && !ws.getClosed() {
		ws.writeChan <- b
	}
}

func (ws *WSession) write(b *buf.Buffer) (err error) {
	src := b.Bytes()
	if encrypt := ws.encrypt; encrypt != nil && b.Encrypt {
		if len(ws.encryptBuf) < len(src) {
			ws.encryptBuf = make([]byte, len(src))
		}
		encrypt(ws.encryptBuf, src)
		_, err = ws.conn.Write(ws.encryptBuf[:len(src)])
	} else {
		_, err = ws.conn.Write(src)
	}
	b.Done()
	return
}

func (ws *WSession) startWriteLoop(startWrite, endWrite chan<- struct{}) {
	startWrite <- struct{}{}
	remain := 0
	for {
		if ws.getClosed() {
			remain = len(ws.writeChan)
			break
		}
		select {
		case b := <-ws.writeChan:
			err := ws.write(b)
			if err != nil {
				ws.setClosed()
			}
		case encrypt := <-ws.encryptChan:
			ws.encrypt = encrypt
		case <-time.After(time.Second):
		}
	}

	for i := 0; i < remain; i++ {
		b := <-ws.writeChan
		err := ws.write(b)
		if err != nil {
			break
		}
	}

	ws.conn.Close()
	log.Error("WSession startWriteLoop end")
	endWrite <- struct{}{}
}

func (ws *WSession) Run(onNewSession, onQuitSession func()) {
	startWrite := make(chan struct{})
	endWrite := make(chan struct{})

	go ws.startWriteLoop(startWrite, endWrite)
	<-startWrite

	if onNewSession != nil {
		onNewSession()
	}

	<-endWrite

	if ws.conn != nil {
		ws.conn.Close()
	}

	if onQuitSession != nil {
		onQuitSession()
	}
}

func (ws *WSession) Quit() {
	ws.setClosed()
}

// NewWSession new a write session
func NewWSession(conn net.Conn, id string, conWriteSize int) *WSession {
	if conWriteSize <= 0 {
		conWriteSize = 4096
	}
	return &WSession{
		conn:        conn,
		id:          id,
		writeQuit:   make(chan struct{}),
		writeChan:   make(chan *buf.Buffer, conWriteSize),
		encryptChan: make(chan EncryptFunc, 1),
		encryptBuf:  make([]byte, 0),
	}
}

// RWSession Read and Write session
type RWSession struct {
	*WSession
	rstream StreamReader
}

// SetDecrypt set decrypt function
func (s *RWSession) SetDecrypt(decrypt DecryptFunc) {
	s.rstream.SetDecrypt(decrypt)
}

func (s *RWSession) startReadLoop(startRead, endRead chan<- struct{}) {
	startRead <- struct{}{}
	for {
		_, err := s.rstream.Read()
		if err != nil {
			s.setClosed()
		}
		if s.getClosed() {
			break
		}
	}
	log.Error("RWSession startReadLoop end")
	endRead <- struct{}{}
}

// Run run session
func (s *RWSession) Run(onNewSession, onQuitSession func()) {
	startRead := make(chan struct{})
	startWrite := make(chan struct{})
	endRead := make(chan struct{})
	endWrite := make(chan struct{})

	go s.startReadLoop(startRead, endRead)
	go s.startWriteLoop(startWrite, endWrite)

	<-startRead
	<-startWrite

	if onNewSession != nil {
		onNewSession()
	}

	<-endRead
	<-endWrite

	if s.conn != nil {
		s.conn.Close()
	}

	if onQuitSession != nil {
		onQuitSession()
	}
}

// NewRWSession new a read and write session
func NewRWSession(conn net.Conn, rstream StreamReader, id string, conWriteSize int) *RWSession {
	s := new(RWSession)
	s.WSession = NewWSession(conn, id, conWriteSize)
	s.rstream = rstream
	return s
}
