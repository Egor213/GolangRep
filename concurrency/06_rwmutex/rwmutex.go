package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	writerBit  = uint32(1) << 31
	readerMask = ^writerBit
)

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		s := atomic.LoadUint32(&rw.state)
		if s&writerBit != 0 {
			futex.Wait(&rw.state, s)
			continue
		}
		if atomic.CompareAndSwapUint32(&rw.state, s, s+1) {
			return
		}
	}
}

func (rw *RWMutex) RUnlock() {
	for {
		s := atomic.LoadUint32(&rw.state)
		rc := s & readerMask
		if rc == 0 {
			panic("rwmutex: RUnlock without RLock")
		}
		if atomic.CompareAndSwapUint32(&rw.state, s, s-1) {
			if rc == 1 && s&writerBit != 0 {
				futex.WakeAll(&rw.state)
			}
			return
		}
	}
}

func (rw *RWMutex) Lock() {
	for {
		s := atomic.LoadUint32(&rw.state)
		if s&writerBit != 0 {
			futex.Wait(&rw.state, s)
			continue
		}
		if atomic.CompareAndSwapUint32(&rw.state, s, s|writerBit) {
			for {
				s2 := atomic.LoadUint32(&rw.state)
				if s2&readerMask == 0 {
					return
				}
				futex.Wait(&rw.state, s2)
			}
		}
	}
}

func (rw *RWMutex) Unlock() {
	for {
		s := atomic.LoadUint32(&rw.state)
		if s&writerBit == 0 {
			panic("rwmutex: Unlock without Lock")
		}
		if atomic.CompareAndSwapUint32(&rw.state, s, s&^writerBit) {
			futex.WakeAll(&rw.state)
			return
		}
	}
}
