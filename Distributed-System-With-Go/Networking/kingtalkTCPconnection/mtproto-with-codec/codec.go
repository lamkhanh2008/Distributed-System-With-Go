package mtproto

type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}
