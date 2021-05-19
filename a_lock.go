package lock_impl

import "sync/atomic"

type ALock struct {
	tail *int32
	flag []bool
	size int
}

func NewALock(capacity int) *ALock {
	l := &ALock{
		tail: new(int32),
		flag: make([]bool, capacity),
		size: capacity,
	}
	l.flag[0] = true
	return l
}

func (a *ALock) Lock() int {
	slot := (int(atomic.AddInt32(a.tail, 1)) - 1) % a.size
	for !a.flag[slot]{}
	return slot
}

func (a *ALock) UnLock(slot int) {
	a.flag[slot] = false
	a.flag[(slot + 1) % a.size] = true
}
