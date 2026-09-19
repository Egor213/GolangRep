package spinlock

import "sync/atomic"

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for !s.TryLock() {
	}
}

func (s *Spinlock) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *Spinlock) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("spinlock: unlock of unlocked lock")
	}
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	for {
		for s.locked.Load() {
		}
		if s.TryLock() {
			return
		}
	}
}

func (s *TTAS) TryLock() bool {
	if s.locked.Load() {
		return false
	}
	return s.locked.CompareAndSwap(false, true)
}

func (s *TTAS) Unlock() {
	if !s.locked.CompareAndSwap(true, false) {
		panic("spinlock: unlock of unlocked lock")
	}
}
