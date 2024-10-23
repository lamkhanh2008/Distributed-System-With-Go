package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

type stateMachine struct {
	db     *sync.Map
	server int
}

type commandKind uint8

const (
	setCommand commandKind = iota
	getCommand
)

type command struct {
	kind  commandKind
	key   string
	value string
}

func (s *stateMachine) Apply(cmd []byte) ([]byte, error) {
	c := decodeCommand(cmd)
	switch c.kind {
	case setCommand:
		s.db.Store(c.key, c.value)

	case getCommand:
		value, ok := s.db.Load(c.value)
		if !ok {
			return nil, fmt.Errorf("Key not found")
		}
		return []byte(value.(string)), nil
	default:
		return nil, fmt.Errorf("Unknown command: %x", cmd)
	}
	return nil, nil
}

func encodeCommand(c command) []byte {
	msg := bytes.NewBuffer(nil)
	err := msg.WriteByte(uint8(c.kind))
	if err != nil {
		panic(err)
	}
	err = binary.Write(msg, binary.LittleEndian, uint64(len(c.key)))
	if err != nil {
		panic(err)
	}
	msg.WriteString(c.key)
	return msg.Bytes()
}

func decodeCommand(msg []byte) command {
	var c command
	c.kind = commandKind(msg[0])
	keyLen := binary.LittleEndian.Uint64(msg[1:9])
	c.key = string(msg[9 : 9+keyLen])
	if c.kind == setCommand {
		valLen := binary.LittleEndian.Uint64(msg[9+keyLen : 9+keyLen+8])
		c.value = string(msg[9+keyLen+8 : 9+keyLen+8+valLen])
	}

	return c
}
