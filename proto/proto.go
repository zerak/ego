package proto

import (
	//"encoding/json"
	"errors"
	"io"

	"github.com/golang/protobuf/proto"
)

const DefaultByteNumForLength = 2

type FactoryFunc func() proto.Message

var (
	protocolFactory map[string]FactoryFunc

	ErrTooShort           = errors.New("too short")
	ErrUnknownMessageName = errors.New("unknown message name")
)

func DecodeLength(buf []byte) int {
	n := 0
	for i, b := range buf {
		n |= int(b) << (uint32(i) << 3)
	}
	return n
}

func EncodeLength(length int, buf []byte) []byte {
	for i := range buf {
		b := length & 0xFF
		buf[i] = byte(b)
		length = length >> (uint32(i+1) << 3)
	}
	return buf
}

func encodeMessageHeader(w io.Writer, v proto.Message, bodySize int) error {
	name := proto.MessageName(v)
	totalSize := DefaultByteNumForLength + (DefaultByteNumForLength + len(name)) + bodySize
	// 写包大小
	if _, err := w.Write(EncodeLength(totalSize, make([]byte, DefaultByteNumForLength))); err != nil {
		return err
	}
	// 写包名大小
	if _, err := w.Write(EncodeLength(len(name), make([]byte, DefaultByteNumForLength))); err != nil {
		return err
	}
	// 写包名
	_, err := io.WriteString(w, name)
	return err
}

func decodeMessageHeader(b []byte) (n int, name string, err error) {
	byteNum := DefaultByteNumForLength
	// 解包大小
	if len(b[n:]) < byteNum {
		err = ErrTooShort
		return
	}
	DecodeLength(b[n : n+byteNum])
	n += byteNum
	// 解包名大小
	if len(b[n:]) < byteNum {
		err = ErrTooShort
		return
	}
	nameLength := DecodeLength(b[n : n+byteNum])
	n += byteNum
	// 解包名
	if len(b[n:]) < nameLength {
		err = ErrTooShort
		return
	}
	name = string(b[n : n+nameLength])
	n += nameLength
	return
}

func Encode(w io.Writer, v proto.Message) error {
	data, err := proto.Marshal(v)
	if err == nil {
		encodeMessageHeader(w, v, len(data))
		_, err = w.Write(data)
	}
	return err
}

func Decode(b []byte) (proto.Message, error) {
	n, name, err := decodeMessageHeader(b)
	if err != nil {
		return nil, err
	}
	fn, ok := protocolFactory[name]
	if !ok {
		return nil, ErrUnknownMessageName
	}
	v := fn()
	err = proto.Unmarshal(b[n:], v)
	return v, err
}
