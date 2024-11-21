package mtproto

import (
	"encoding/binary"
	"fmt"
	"io"
)

type MTProtoIntermediateCodec struct {
	conn *BufferedConn
}

func NewMTProtoIntermediateCodec(conn *BufferedConn) *MTProtoIntermediateCodec {
	return &MTProtoIntermediateCodec{
		conn: conn,
	}
}
func (c *MTProtoIntermediateCodec) Receive() (interface{}, error) {
	var size int
	var n int
	var err error

	b := make([]byte, 4)
	n, err = io.ReadFull(c.conn, b)
	if err != nil {
		return nil, err
	}

	size = int(binary.LittleEndian.Uint32(b))

	// glog.V(3).Info("first_byte: ", hex.EncodeToString(b[:1]))
	// needAck := bool(b[0] >> 7 == 1)
	// _ = needAck

	//b[0] = b[0] & 0x7f
	//// glog.V(3).Info("first_byte2: ", hex.EncodeToString(b[:1]))
	//
	//if b[0] < 0x7f {
	//	size = int(b[0]) << 2
	//	glog.V(3).Info("size1: ", size)
	//	if size == 0 {
	//		return nil, nil
	//	}
	//} else {
	//	glog.V(3).Info("first_byte2: ", hex.EncodeToString(b[:1]))
	//	b2 := make([]byte, 3)
	//	n, err = io.ReadFull(c.conn, b2)
	//	if err != nil {
	//		return nil, err
	//	}
	//	size = (int(b2[0]) | int(b2[1])<<8 | int(b2[2])<<16) << 2
	//	glog.V(3).Info("size2: ", size)
	//}

	left := size
	buf := make([]byte, size)
	for left > 0 {
		n, err = io.ReadFull(c.conn, buf[size-left:])
		if err != nil {
			fmt.Println("ReadFull2 error: ", err)
			return nil, err
		}
		left -= n
	}
	//if size > 10240 {
	//	glog.V(3).Info("ReadFull2: ", hex.EncodeToString(buf[:256]))
	//}

	// TODO(@NhokCrazy199): process report ack and quickack
	// Truncating the QuickAck message, the client has a problem
	if size == 4 {
		fmt.Println("Server response error: ", int32(binary.LittleEndian.Uint32(buf)))
		// return nil, fmt.Errorf("Recv QuickAckMessage, ignore!!!!") //  connId: ", c.stream, ", by client ", m.RemoteAddr())
		return nil, nil
	}

	authKeyId := int64(binary.LittleEndian.Uint64(buf))
	message := NewMTPRawMessage(authKeyId, 0, TRANSPORT_TCP)
	message.Decode(buf)
	return message, nil
}

func (c *MTProtoIntermediateCodec) Send(msg interface{}) error {
	message, ok := msg.(*MTPRawMessage)
	if !ok {
		err := fmt.Errorf("msg type error, only MTPRawMessage, msg: {%v}", msg)
		fmt.Println(err)
		return err
	}

	b := message.Encode()

	sb := make([]byte, 4)
	// minus padding
	size := len(b)

	//if size < 127 {
	//	sb = []byte{byte(size)}
	//} else {
	binary.LittleEndian.PutUint32(sb, uint32(size))
	//}

	b = append(sb, b...)
	_, err := c.conn.Write(b)

	if err != nil {
		fmt.Println("Send msg error: %s", err)
	}

	return err
}

func (c *MTProtoIntermediateCodec) Close() error {
	return c.conn.Close()
}
