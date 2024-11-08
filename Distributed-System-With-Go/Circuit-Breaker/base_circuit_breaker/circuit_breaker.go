package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrRefuse = errors.New("request refused. the circuit is open")
)

type CircuitBreaker interface {
	Execute(func() (interface{}, error)) (interface{}, error)
	State() string
}

type Policy int
type State string

const (
	MaxFails Policy = iota
	MaxConsecutiveFails
)

const (
	open     State = "open"
	halfOpen State = "half-open"
	close    State = "closed"
)

type ExtraOptions struct {
	Policy              Policy
	MaxFails            *uint64
	MaxConsecutiveFails *uint64
	OpenInterval        *time.Duration
}

type circuitbreaker struct {
	policy              Policy
	maxFails            uint64
	maxConsecutiveFails uint64
	openInterval        time.Duration
	fails               uint64
	state               State
	openChannel         chan struct{}
	mutex               sync.Mutex
}

func ToPointer[T any](l T) *T {
	return &l
}

func New(opts ...ExtraOptions) CircuitBreaker {
	var opt ExtraOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.MaxFails == nil {
		opt.MaxFails = ToPointer(uint64(5))
	}

	if opt.MaxConsecutiveFails == nil {
		opt.MaxConsecutiveFails = ToPointer(uint64(5))
	}

	if opt.OpenInterval == nil {
		opt.OpenInterval = ToPointer(5 * time.Second)
	}

	cb := &circuitbreaker{
		policy:              opt.Policy,
		maxFails:            *opt.MaxFails,
		maxConsecutiveFails: *opt.MaxConsecutiveFails,
		openInterval:        *opt.OpenInterval,
		openChannel:         make(chan struct{}),
	}

	go cb.openWatcher()

	return cb
}

func (cb *circuitbreaker) State() string {
	return string(cb.state)
}

func (cb *circuitbreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	err := cb.doPreRequest()
	if err != nil {
		return nil, err
	}

	res, err := req()
	err = cb.doPostRequest(err)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (cb *circuitbreaker) doPreRequest() error {
	if cb.state == open {
		return ErrRefuse
	}
	return nil
}

func (cb *circuitbreaker) doPostRequest(err error) error {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err == nil {
		if cb.policy == MaxConsecutiveFails {
			cb.fails = 0
		}
		cb.state = close
		return nil
	}
	if cb.state == halfOpen {
		//nếu request process fails thì chuyển sang open, còn success thì sang closed
		cb.state = open
		cb.openChannel <- struct{}{}
		return err
	}

	cb.fails += 1
	if cb.failsExcceededThreshod() {
		cb.state = open
		cb.openChannel <- struct{}{}
	}

	return err
}

func (cb *circuitbreaker) failsExcceededThreshod() bool {
	switch cb.policy {
	case MaxConsecutiveFails:
		return cb.fails >= cb.maxConsecutiveFails
	case MaxFails:
		return cb.fails >= cb.maxFails
	default:
		return false
	}
}

func (cb *circuitbreaker) openWatcher() {
	for range cb.openChannel {
		time.Sleep(cb.openInterval)
		cb.mutex.Lock()
		cb.state = halfOpen
		cb.fails = 0
		cb.mutex.Unlock()
	}
}
