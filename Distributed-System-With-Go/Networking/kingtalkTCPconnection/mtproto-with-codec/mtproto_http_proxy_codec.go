package mtproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type MTProtoHttpProxyCodec struct {
	conn net.Conn
}

func NewMTProtoHttpProxyCodec(conn net.Conn) *MTProtoHttpProxyCodec {
	// return conn.tcpConn.SetReadDeadline(time.Now().Add(tcpHeartbeat * 2))
	conn.SetReadDeadline(time.Now().Add(time.Second * 60))
	// .(*ktcp.BufferedConn).Conn.(*net.TCPConn).SetReadDeadline(time.Now().Add())
	return &MTProtoHttpProxyCodec{
		conn: conn,
	}
}

func (c *MTProtoHttpProxyCodec) Receive() (interface{}, error) {
	req, err := http.ReadRequest(c.conn.(*BufferedConn).BufioReader())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(body) < 8 {
		err = fmt.Errorf("not enough uint64 len error - %d", len(body))
		fmt.Println(err)
		return nil, err
	}
	authKeyId := int64(binary.LittleEndian.Uint64(body))
	msg := NewMTPRawMessage(authKeyId, 0, TRANSPORT_HTTP)
	err = msg.Decode(body)
	fmt.Println("--------------Done MTP receiver: %v", msg)
	if err != nil {
		fmt.Println(err)
		// conn.Close()
		return nil, err
	}

	return msg, nil
}

func (c *MTProtoHttpProxyCodec) Send(msg interface{}) error {
	message, ok := msg.(*MTPRawMessage)
	if !ok {
		err := fmt.Errorf("msg type error, only MTPRawMessage, msg: {%v}", msg)
		fmt.Println(err)
		// conn.Close()
		return err
	}
	b := message.Encode()
	rsp := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Request:    &http.Request{Method: "POST"},
		Header: http.Header{
			"Access-Control-Allow-Headers": {"origin, content-type"},
			"Access-Control-Allow-Methods": {"POST, OPTIONS"},
			"Access-Control-Allow-Origin":  {"*"},
			"Access-Control-Max-Age":       {"1728000"},
			"Cache-control":                {"no-store"},
			"Connection":                   {"keep-alive"},
			"Content-type":                 {"application/octet-stream"},
			"Pragma":                       {"no-cache"},
			"Strict-Transport-Security":    {"max-age=15768000"},
		},
		ContentLength: int64(len(b)),
		Body:          ioutil.NopCloser(bytes.NewReader(b)),
		Close:         false,
	}
	err := rsp.Write(c.conn)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (c *MTProtoHttpProxyCodec) Close() error {
	return c.conn.Close()
}
