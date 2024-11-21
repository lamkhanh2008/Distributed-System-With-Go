package mtproto

import (
	"fmt"
	"io"
)

type Protocol interface {
	NewCodec(rw io.ReadWriter) (Codec, error)
}

var (
	protocolRegister = make(map[string]Protocol)
)

func RegisterProtocol(name string, protocol Protocol) {
	fmt.Println(name, protocol)
	protocolRegister[name] = protocol
}

func NewCodecByName(name string, rw io.ReadWriter) (Codec, error) {
	protocol, ok := protocolRegister[name]
	if !ok {
		return nil, fmt.Errorf("not found protocol name: %s", name)
	}
	return protocol.NewCodec(rw)
}
