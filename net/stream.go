package net

import (
	"io"
	"net"
	"sync"
	"time"

	"ego/proto"
)

type StreamReader interface {
	Conn() net.Conn
	Read() (n int, err error)
	SetDecrypt(decrypt DecryptFunc)
	SetTimeout(d time.Duration)
}

type ReadStream struct {
	conn             net.Conn
	timeout          time.Duration
	bufForLength     []byte
	buf              []byte
	byteNumForLength int
	packetHandler    PacketHandler
	decodeLengthFunc func([]byte) int

	decryptLocker sync.RWMutex
	decrypt       DecryptFunc
}

func (r *ReadStream) Conn() net.Conn { return r.conn }
func (r *ReadStream) Read() (int, error) {
	total := 0
	// 2 bytes size represent packet size
	if r.timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.timeout))
	}
	n, err := r.conn.Read(r.buf[:r.byteNumForLength])
	total += n
	if err != nil {
		return total, err
	}

	// get current decrypt
	r.decryptLocker.RLock()
	decrypter := r.decrypt
	r.decryptLocker.RUnlock()

	// encrypt packet size
	if decrypter != nil {
		decrypter(r.buf[:r.byteNumForLength], r.buf[:r.byteNumForLength])
	}
	nextLength := r.decodeLengthFunc(r.buf[:r.byteNumForLength])
	readedSize := r.byteNumForLength
	if len(r.buf) < nextLength {
		r.buf = make([]byte, nextLength)
	}

	// according to packet size parse packet data
	for readedSize < nextLength {
		n, err := r.conn.Read(r.buf[readedSize:nextLength])
		total += n
		readedSize += n
		if err != nil {
			return total, err
		}
	}
	if decrypter != nil {
		decrypter(r.buf[r.byteNumForLength:nextLength], r.buf[r.byteNumForLength:nextLength])
	}
	r.packetHandler.OnPacket(r.buf[:nextLength])
	return total, nil
}
func (r *ReadStream) SetDecrypt(decrypt DecryptFunc) {
	r.decryptLocker.Lock()
	defer r.decryptLocker.Unlock()
	r.decrypt = decrypt
}
func (r *ReadStream) SetTimeout(d time.Duration) { r.timeout = d }
func (r *ReadStream) SetByteNumForLength(n int) {
	r.byteNumForLength = n
	if len(r.buf) < n {
		r.buf = make([]byte, n)
	}
}
func (r *ReadStream) SetDecodeLengthFunc(decodeLengthFunc func([]byte) int) {
	r.decodeLengthFunc = decodeLengthFunc
}

// 通用的流读取
func NewReadStream(conn net.Conn, onPacket PacketHandler) *ReadStream {
	return &ReadStream{
		conn:             conn,
		buf:              make([]byte, 4096),
		byteNumForLength: proto.DefaultByteNumForLength,
		packetHandler:    onPacket,
		decodeLengthFunc: proto.DecodeLength,
	}
}

// UDP 包读取
type UDPReadStream struct {
	conn          *net.UDPConn
	packetHandler PacketHandler
	buf           []byte
	timeout       time.Duration
}

func (r *UDPReadStream) Conn() net.Conn { return r.conn }
func (r *UDPReadStream) Read(encryptFunc EncryptFunc) (int, error) {
	total := 0
	if r.timeout > 0 {
		r.conn.SetReadDeadline(time.Now().Add(r.timeout))
	}
	n, _, err := r.conn.ReadFromUDP(r.buf)
	total += n
	if err != nil && err != io.EOF {
		return total, err
	}
	if n > 0 {
		r.packetHandler.OnPacket(r.buf[:n])
	}
	return total, nil
}
func (r *UDPReadStream) SetDecrypter(EncryptFunc) {
	panic("UDPReadStream unsupport decrypter")
}
func (r *UDPReadStream) SetTimeout(d time.Duration) { r.timeout = d }

func NewUDPReadStream(conn *net.UDPConn, onPacket PacketHandler) *UDPReadStream {
	return &UDPReadStream{
		conn:          conn,
		packetHandler: onPacket,
		buf:           make([]byte, 8192),
	}
}
