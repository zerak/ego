package net

type PacketHandler interface {
	OnPacket([]byte)
}

type DefaultPacketHandler struct {
}

func (m *DefaultPacketHandler) OnPacket(b []byte) {

}
func NewDefaultPacketHandler() *DefaultPacketHandler {
	return &DefaultPacketHandler{}
}
