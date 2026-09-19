package waitgroup

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	newCount := atomic.AddUint32(&wg.count, uint32(delta))

	if int32(newCount) < 0 {
		panic("waitgroup: negative counter")
	}

	if newCount == 0 {
		futex.WakeAll(&wg.count)
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	for {
		v := atomic.LoadUint32(&wg.count)
		if v == 0 {
			return
		}

		futex.Wait(&wg.count, v)
	}
}
