package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

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

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go increment(i, m, &wg)
	}
	wg.Wait()

}

var (
	x int
	m = &TTasLock{
		State: new(int32),
	}
)

func increment(i int, m *TTasLock, wg *sync.WaitGroup) {
	for {
		m.Lock()
		if x >= 200 {
			m.UnLock()
			wg.Done()
			return
		}
		fmt.Println(x)
		x++
		m.UnLock()
	}
}
