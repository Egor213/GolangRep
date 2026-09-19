package barrier

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Barrier struct {
	need    uint32
	arrived uint32
	round   uint32
}

func New(n int) *Barrier {
	if n <= 0 {
		panic("barrier: n must be positive")
	}
	return &Barrier{need: uint32(n)}
}

func (b *Barrier) Wait() {
	r := atomic.LoadUint32(&b.round)
	n := atomic.AddUint32(&b.arrived, 1)

	if n == b.need {
		atomic.StoreUint32(&b.arrived, 0)
		atomic.AddUint32(&b.round, 1)
		futex.WakeAll(&b.round)
		return
	}

	for atomic.LoadUint32(&b.round) == r {
		futex.Wait(&b.round, r)
	}
}
