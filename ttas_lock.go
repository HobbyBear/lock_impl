package lock_impl

import "sync/atomic"

type TTasLock struct {
	State *int32
}

func (t *TTasLock) Lock() {
	for true {
		for atomic.LoadInt32(t.State) == 1 {
		}
		if atomic.CompareAndSwapInt32(t.State, 0, 1) {
			return
		}
	}
}

func (t *TTasLock) UnLock() {
	atomic.StoreInt32(t.State, 0)
}
