package mutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	if atomic.CompareAndSwapUint32(&m.state, free, held) {
		return
	}

	for range 50 {
		if atomic.LoadUint32(&m.state) == free {
			if atomic.CompareAndSwapUint32(&m.state, free, held) {
				return
			}
		}
	}

	for {
		if atomic.CompareAndSwapUint32(&m.state, free, contended) {
			return
		}

		atomic.CompareAndSwapUint32(&m.state, held, contended)

		futex.Wait(&m.state, contended)
	}

}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	for {
		state := atomic.LoadUint32(&m.state)

		switch state {
		case free:
			panic("mutex: unlock of unlocked mutex")
		case held:
			if atomic.CompareAndSwapUint32(&m.state, held, free) {
				return
			}
		case contended:
			if atomic.CompareAndSwapUint32(&m.state, contended, free) {
				futex.Wake(&m.state)
				return
			}
		}
	}

}
