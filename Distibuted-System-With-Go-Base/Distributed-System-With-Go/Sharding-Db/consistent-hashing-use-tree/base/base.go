package consistent

import (
	"errors"
	"sync"
)

type uints []uint32

func (x uints) Len() int { return len(x) }

func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var ErrEmptyCircle = errors.New("empty circle")

type Consistent struct {
	circle           map[uint32]string
	members          map[string]bool
	sortedHashes     uints
	NumberOfReplicas int
	count            int64
	scratch          [64]byte
	UseFnv           bool
	sync.RWMutex
}
