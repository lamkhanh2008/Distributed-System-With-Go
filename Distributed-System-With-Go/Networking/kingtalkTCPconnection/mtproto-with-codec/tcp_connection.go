package mtproto

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// ErrorConnectionClosed = errors.New("Connection Closed")
var ErrorConnectionClosed = errors.New("Connection Closed")

// ErrorConnectionBlocked = errors.New("Connection Blocked")
var ErrorConnectionBlocked = errors.New("Connection Blocked")

var globalConnectionID uint64

type TCPConnection struct {
	name          string
	conn          net.Conn
	id            uint64
	codec         Codec
	sendChan      chan interface{}
	recvMutex     sync.Mutex
	sendMutex     sync.RWMutex
	closeFlag     int32
	closeChan     chan int
	closeMutex    sync.Mutex
	closeCallback closeCallback
	Context       interface{}
}
type ClearSendChan interface {
	ClearSendChan(<-chan interface{})
}

type closeCallback interface {
	// func(Connection)
	OnConnectionClosed(Connection)
}

type Connection interface {
	GetConnID() uint64
	IsClosed() bool
	Close() error
	Codec() Codec
	Receive() (interface{}, error)
	Send(msg interface{}) error
}

func NewTCPConnection(name string, conn net.Conn, sendChanSize int, codec Codec, cb closeCallback) *TCPConnection {
	if globalConnectionID >= 0xfffffffffffffff {
		atomic.StoreUint64(&globalConnectionID, 0)
	}

	conn2 := &TCPConnection{
		name:      name,
		conn:      conn,
		codec:     codec,
		closeChan: make(chan int),
		// id:        atomic.AddUint64(&globalConnectionID, 1),
		id:            genSessionID(name, atomic.AddUint64(&globalConnectionID, 1)),
		closeCallback: cb,
	}

	if sendChanSize > 0 {
		conn2.sendChan = make(chan interface{}, sendChanSize)
		go conn2.sendLoop()
	}
	return conn2
}

func genSessionID(name string, id uint64) uint64 {
	var sid = id
	if name == "frontend443" {
		// sid = sid | 0 << 60
	} else if name == "frontend80" {
		sid = sid | 1<<60
	} else if name == "frontend5222" {
		sid = sid | 2<<60
	}
	return sid
}

func (c *TCPConnection) String() string {
	return fmt.Sprintf("{connID: %d@%s-(%s->%s)}", c.id, c.name, c.conn.LocalAddr(), c.conn.RemoteAddr())
}

// LoadAddr func
func (c *TCPConnection) LoadAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr func
func (c *TCPConnection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// Name func
func (c *TCPConnection) Name() string {
	return c.name
}

// GetConnID func
func (c *TCPConnection) GetConnID() uint64 {
	return c.id
}

// GetNetConn func
func (c *TCPConnection) GetNetConn() net.Conn {
	return c.conn
}

// GetContext func
func (c *TCPConnection) GetContext() interface{} {
	return c.Context
}

// SetContext func
func (c *TCPConnection) SetContext(context interface{}) {
	c.Context = context
}

// IsClosed func
func (c *TCPConnection) IsClosed() bool {
	return atomic.LoadInt32(&c.closeFlag) == 1
}

// Close func
func (c *TCPConnection) Close() error {
	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
		if c.closeCallback != nil {
			c.closeCallback.OnConnectionClosed(c)
		}

		close(c.closeChan)

		if c.sendChan != nil {
			c.sendMutex.Lock()
			close(c.sendChan)
			if clear, ok := c.codec.(ClearSendChan); ok {
				clear.ClearSendChan(c.sendChan)
			}
			c.sendMutex.Unlock()
		}

		err := c.codec.Close()
		return err
	}
	return ErrorConnectionClosed
}

// Codec func
func (c *TCPConnection) Codec() Codec {
	return c.codec
}

// Receive func
func (c *TCPConnection) Receive() (interface{}, error) {
	c.recvMutex.Lock()
	defer c.recvMutex.Unlock()

	msg, err := c.codec.Receive()
	if err != nil {
		c.Close()
	}
	fmt.Println("------------------Done receive %v", msg)
	return msg, err
}

// sendLoop func
func (c *TCPConnection) sendLoop() {
	defer c.Close()
	for {
		select {
		case msg, ok := <-c.sendChan:
			if !ok || c.codec.Send(msg) != nil {
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

// Send func
func (c *TCPConnection) Send(msg interface{}) error {
	if c.sendChan == nil {
		if c.IsClosed() {
			return ErrorConnectionClosed
		}

		c.sendMutex.Lock()
		defer c.sendMutex.Unlock()

		err := c.codec.Send(msg)
		if err != nil {
			c.Close()
		}
		return err
	}

	c.sendMutex.RLock()
	if c.IsClosed() {
		c.sendMutex.RUnlock()
		return ErrorConnectionClosed
	}

	select {
	case c.sendChan <- msg:
		c.sendMutex.RUnlock()
		return nil
	default:
		c.sendMutex.RUnlock()
		c.Close()
		return ErrorConnectionBlocked
	}
}
