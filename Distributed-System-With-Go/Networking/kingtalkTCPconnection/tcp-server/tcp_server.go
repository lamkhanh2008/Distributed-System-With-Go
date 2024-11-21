package tcp

import (
	"fmt"
	"io"
	"kingtalk_tcp/mtproto-with-codec"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const maxConcurrentConnection = 10000000

type TCPServer struct {
	connectionManager *ConnectionManager
	listener          net.Listener
	serverName        string
	protoName         string
	sendChanSize      int
	callback          TCPConnectionCallback
	running           bool
	sem               chan struct{}
	releaseOnce       sync.Once
}

type Connection interface {
	GetConnID() uint64
	IsClosed() bool
	Close() error
	Codec() mtproto.Codec
	Receive() (interface{}, error)
	Send(msg interface{}) error
}

// closeCallback interface
type closeCallback interface {
	// func(Connection)
	OnConnectionClosed(Connection)
}

type TCPConnectionCallback interface {
	OnNewConnection(conn *TCPConnection)
	OnConnectionDataArrived(c *TCPConnection, msg interface{}) error
	OnConnectionClosed(c *TCPConnection)
}
type TCPServerArgs struct {
	Listener                net.Listener
	ServerName              string
	ProtoName               string
	SendChanSize            int
	ConnectionCallback      TCPConnectionCallback
	MaxConcurrentConnection int
}

func NewTCPServer(args TCPServerArgs) *TCPServer {
	if args.MaxConcurrentConnection < 1 {
		args.MaxConcurrentConnection = maxConcurrentConnection
	}
	return &TCPServer{
		connectionManager: NewConnectionManager(),
		listener:          args.Listener,
		serverName:        args.ServerName,
		protoName:         args.ProtoName,
		sendChanSize:      args.SendChanSize,
		callback:          args.ConnectionCallback,
		running:           false,
		sem:               make(chan struct{}, args.MaxConcurrentConnection),
	}
}

func (s *TCPServer) IsRunning() bool {
	return s.running
}

func (s *TCPServer) Serve() {
	if s.running {
		return
	}
	s.running = true
	s.acquire()
	for {
		conn, err := Accept(s.listener)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn2 := mtproto.NewBufferedConn(conn)
		codec, err := mtproto.NewCodecByName(s.protoName, conn2)
		if err != nil {
			fmt.Println(err)
			s.running = false
			conn.Close()
			return
		}
		tcpConn := NewTCPConnection(s.serverName, conn2, s.sendChanSize, codec, s)

		go s.establishTCPConnection(tcpConn)

	}
}

func (s *TCPServer) acquire() { s.sem <- struct{}{} }

// release func
func (s *TCPServer) release() { <-s.sem }

func Accept(listener net.Listener) (net.Conn, error) {
	var tempDelay time.Duration
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil, io.EOF
			}
			return nil, err
		}
		return conn, nil
	}
}

func (s *TCPServer) establishTCPConnection(conn *TCPConnection) {
	// glog.V(3).Info("establishTcpConnection...")
	defer func() {
		//
		if err := recover(); err != nil {
			fmt.Println("tcp_server handle panic: %v\n%s", err, debug.Stack())
			conn.Close()
		}
	}()

	s.onNewConnection(conn)
	fmt.Printf("Connected")
	for {
		conn.conn.SetReadDeadline(time.Now().Add(time.Minute * 6))
		msg, err := conn.Receive()
		if err != nil {
			fmt.Println("conn {%v} recv error: %v", conn, err)
			return
		}

		if msg == nil {
			fmt.Println("recv a nil msg: %v", conn)
			// Do you need to close it?
			continue
		}

		// go func() {
		// 	newspan := opentracing.StartSpan(
		// 		"establishTCPConnection",
		// 	)
		// 	newspan.SetTag("lalalala2", "lolololo2")

		// 	_ = newspan
		// 	defer newspan.Finish()
		// }()

		if s.callback != nil {
			if err := s.callback.OnConnectionDataArrived(conn, msg); err != nil {
				// TODO: Do you need to close it?
				fmt.Println("OnConnectionDataArrived err: %v for conn: %v", err, conn)
			}
		} else {
			fmt.Println("s.callback is nil for conn: %v", conn)
		}
	}
}

func (s *TCPServer) onNewConnection(conn *TCPConnection) {
	if s.connectionManager != nil {
		s.connectionManager.PutConnection(conn)
	}

	if s.callback != nil {
		s.callback.OnNewConnection(conn)
	}
}

func (s *TCPServer) OnConnectionClosed(conn Connection) {
	s.onConnectionClosed(conn.(*TCPConnection))
}

func (s *TCPServer) onConnectionClosed(conn *TCPConnection) {
	if s.connectionManager != nil {
		s.connectionManager.delConnection(conn)
	}

	if s.callback != nil {
		s.callback.OnConnectionClosed(conn)
	}
}
