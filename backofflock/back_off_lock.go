package main

import (
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

type BackOffLock struct {
	State   *int32
	backOff *backOff
}

func NewBackOffLock() *BackOffLock {
	return &BackOffLock{
		State:   new(int32),
		backOff: newBackOff(10, 2),
	}
}

func (t *BackOffLock) Lock() {
	for true {
		for atomic.LoadInt32(t.State) == 1 {
		}
		if atomic.CompareAndSwapInt32(t.State, 0, 1) {
			return
		} else {
			t.backOff.backoff()
		}
	}
}

func (t *BackOffLock) UnLock() {
	atomic.StoreInt32(t.State, 0)
}

type backOff struct {
	MinDelay int
	MaxDelay int
	limit    int
}

func newBackOff(min, max int) *backOff {
	return &backOff{
		MinDelay: min,
		MaxDelay: max,
		limit:    min,
	}
}

func (b *backOff) backoff() {
	delay := rand.Intn(b.limit)
	b.limit = int(math.Min(float64(b.MaxDelay), float64(2*b.limit)))
	time.Sleep(time.Duration(delay) * time.Second)
}

func main() {

}
