package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	if n < 0 {
		panic("semaphore: negative permits")
	}
	return &Semaphore{
		permits: uint32(n),
	}
}

func (s *Semaphore) Acquire() {
	for {

		p := atomic.LoadUint32(&s.permits)
		if p == 0 {
			futex.Wait(&s.permits, 0)
		} else {
			if atomic.CompareAndSwapUint32(&s.permits, p, p-1) {
				break
			}
		}

	}
}

func (s *Semaphore) TryAcquire() bool {
	p := atomic.LoadUint32(&s.permits)
	if p > 0 {
		return atomic.CompareAndSwapUint32(&s.permits, p, p-1)
	}
	return false
}

func (s *Semaphore) Release() {
	atomic.AddUint32(&s.permits, 1)
	futex.Wake(&s.permits)
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
