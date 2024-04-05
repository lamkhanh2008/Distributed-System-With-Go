package session_client

import (
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/zrpc"
)

type Session struct {
	dispatcher  *hash.ConsistentHash
	errNotFound error
	sessions    map[string]SessionClient
}

func NewSession(c zrpc.RpcClientConf) *Session {
	sess := &Session{
		dispatcher:  hash.NewConsistentHash(),
		errNotFound: ErrSessionNotFound,
		sessions:    make(map[string]SessionClient),
	}
	sess.watch(c)

	return sess
}
func (see *Session) watch(c zrpc.RpcClientConf) {
	sub, _ := discov.NewSubscriber(c.Etcd.Hosts, c.Etcd.Key)
	update := func() {
		values := sub.Values()
		if len(values) == 0 {
			return
		}
		fmt.Print(values)

		for _, v := range values {
			c.Target = string(v)

			// cli, err := zrpc.NewClient(c)
			// if err != nil {
			// 	fmt.Printf("watchComet NewClient(%+v) error(%v)", values, err)
			// 	return
			// }
			// sessionCli := NewSessionClient(cli)
			// see.dispatcher.Add(sessionCli)
			see.dispatcher.Add(c.Target)
		}
		// cli, err := zrpc.NewClient(c)
		// if err != nil {
		// 	fmt.Printf("watchComet NewClient(%+v) error(%v)", values, err)
		// 	return
		// }
		// sessionCli := NewSessionClient(cli)
		// see.dispatcher.Add(sessionCli)
	}

	sub.AddListener(update)
	update()
}

var (
	ErrSessionNotFound = errors.New("not found session")
)

func (see *Session) GetSessionClient(key string, c zrpc.RpcClientConf) (SessionClient, error) {
	val, _ := see.dispatcher.Get(key)
	c.Target = val.(string)
	cli, err := zrpc.NewClient(c)
	if err != nil {
		fmt.Printf("watchComet NewClient(%+v) error(%v)", err)
		return nil, err
	}
	sessionCli := NewSessionClient(cli)
	return sessionCli, nil
}
