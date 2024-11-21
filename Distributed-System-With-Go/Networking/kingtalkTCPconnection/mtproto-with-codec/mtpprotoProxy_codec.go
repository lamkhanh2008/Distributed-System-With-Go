package mtproto

import (
	"fmt"
	"io"
	"net"
)

const (
	// Tcp Transport
	MTPROTO_ABRIDGED_FLAG     = 0xef
	MTPROTO_INTERMEDIATE_FLAG = 0xeeeeeeee

	// Http Transport
	HTTP_HEAD_FLAG   = 0x44414548
	HTTP_POST_FLAG   = 0x54534f50
	HTTP_GET_FLAG    = 0x20544547
	HTTP_OPTION_FLAG = 0x4954504f

	VAL2_FLAG = 0x00000000
)

const (
	TRANSPORT_TCP  = 1 // TCP
	TRANSPORT_HTTP = 2 // HTTP
	TRANSPORT_UDP  = 3 // UDP, @NhokCrazy199: No UDP-enabled clients are found
)

func init() {
	RegisterProtocol("mtproto", NewMTPProtoProxy())
}

type MTProtoProxy struct {
}

func NewMTPProtoProxy() *MTProtoProxy {
	return &MTProtoProxy{}
}

func (m *MTProtoProxy) NewCodec(rw io.ReadWriter) (Codec, error) {
	codec := &MTProtoProxyCodec{
		codecType: TRANSPORT_TCP,
		conn:      rw.(net.Conn),
		State:     STATE_CONNECTED,
		proto:     m,
	}
	return codec, nil
}

const (
	STATE_CONNECTED = iota
	// STATE_FIRST_BYTE		//
	// STATE_FIRST_INT32		//
	// STATE_FIRST_64BYTES		//
	STATE_DATA //
)

type MTProtoProxyCodec struct {
	codecType int // codec type
	conn      net.Conn
	State     int
	codec     Codec
	proto     *MTProtoProxy
}

func (c *MTProtoProxyCodec) peekCodec() error {
	//if c.State == STATE_DATA {
	//	return nil
	//}
	conn, _ := c.conn.(*BufferedConn)
	// var b_0_1 = make([]byte, 1)
	b_0_1, err := conn.Peek(1)

	// _, err := io.ReadFull(c.conn, b_0_1)
	if err != nil {
		fmt.Println("MTProtoProxyCodec - read b_0_1 error: %v", err)
		return err
	}

	// if b_0_1[0] == MTPROTO_ABRIDGED_FLAG {
	// 	glog.V(3).Info("mtproto abridged version.")
	// 	c.codec = NewMTProtoAbridgedCodec(conn)
	// 	conn.Discard(1)
	// 	return nil
	// }

	// not abridged version, we'll lookup codec!
	// b_1_3 = make([]byte, 3)
	b_1_3, err := conn.Peek(4)
	if err != nil {
		fmt.Println("MTProtoProxyCodec - read b_1_3 error: %v", err)
		return err
	}

	b_1_3 = b_1_3[1:4]
	// first uint32
	val := (uint32(b_1_3[2]) << 24) | (uint32(b_1_3[1]) << 16) | (uint32(b_1_3[0]) << 8) | (uint32(b_0_1[0]))
	if val == HTTP_HEAD_FLAG || val == HTTP_POST_FLAG || val == HTTP_GET_FLAG || val == HTTP_OPTION_FLAG {
		// http
		fmt.Println("mtproto http.")

		// conn2 := NewMTProtoHttpProxyConn(conn)
		// c.conn = conn2
		c.codecType = TRANSPORT_HTTP
		c.codec = NewMTProtoHttpProxyCodec(c.conn)

		// c.proto.httpListener.acceptChan <- conn2
		return nil
	}

	// an intermediate version
	if val == MTPROTO_INTERMEDIATE_FLAG {
		//glog.V(3).Infof("MTProtoProxyCodec - mtproto intermediate version, impl in the future!!")
		//return nil, errors.New("mtproto intermediate version not impl!!")
		fmt.Println("mtproto intermediate version.")
		c.codec = NewMTProtoIntermediateCodec(conn)
		conn.Discard(4)
		return nil
	}

	// recv 4~64 bytes
	// var b_4_60 = make([]byte, 60)
	// b_4_60, err := conn.Peek(64)
	// // io.ReadFull(c.conn, b_4_60)
	// if err != nil {
	// 	glog.V(1).Infof("MTProtoProxyCodec - read b_4_60 error: %v", err)
	// 	return err
	// }
	// b_4_60 = b_4_60[4:64]
	// val2 := (uint32(b_4_60[3]) << 24) | (uint32(b_4_60[2]) << 16) | (uint32(b_4_60[1]) << 8) | (uint32(b_4_60[0]))
	// if val2 == VAL2_FLAG {
	// 	glog.V(3).Info("mtproto full version.")
	// 	c.codec = NewMTProtoFullCodec(conn)
	// 	return nil
	// }

	// var tmp [64]byte
	// // generate decrypt_key
	// for i := 0; i < 48; i++ {
	// 	tmp[i] = b_4_60[51-i]
	// }

	// e, err := crypto.NewAesCTR128Encrypt(tmp[:32], tmp[32:48])
	// if err != nil {
	// 	// glog.V(1).Info("NewAesCTR128Encrypt error: %s", err)
	// 	return err
	// }

	// d, err := crypto.NewAesCTR128Encrypt(b_4_60[4:36], b_4_60[36:52])
	// if err != nil {
	// 	glog.V(1).Infof("NewAesCTR128Encrypt error: %s", err)
	// 	return err
	// }

	// d.Encrypt(b_0_1)
	// d.Encrypt(b_1_3)
	// d.Encrypt(b_4_60)

	// if b_4_60[52] != 0xef && b_4_60[53] != 0xef && b_4_60[54] != 0xef && b_4_60[55] != 0xef {
	// 	glog.V(1).Infof("MTProtoProxyCodec - first 56~59 byte != 0xef")
	// 	return errors.New("mtproto buf[56:60]'s byte != 0xef!!")
	// }

	// glog.V(3).Info("first_bytes_64: ", hex.EncodeToString(b_0_1), hex.EncodeToString(b_1_3), hex.EncodeToString(b_4_60))
	// c.codec = NewMTProtoAppCodec(conn, d, e)
	// conn.Discard(64)

	return nil
}

func (c *MTProtoProxyCodec) Receive() (interface{}, error) {
	if c.codec == nil {
		err := c.peekCodec()
		if err != nil {
			return nil, err
		}
	}
	res, err := c.codec.Receive()
	if err != nil {
		fmt.Println("MTProtoProxyCodec::Receive Codec: %T - Error %T - %+v - Res: %T -%+v", c.codec, err, err, res, res)
	}
	
	return res, err
}

func (c *MTProtoProxyCodec) Send(msg interface{}) error {
	if c.codec != nil {
		return c.codec.Send(msg)
	}
	return fmt.Errorf("codec is nil")
}

func (c *MTProtoProxyCodec) Close() error {
	if c.codec != nil {
		return c.codec.Close()
	} else {
		return nil
	}
}
