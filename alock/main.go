package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

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
	m = NewALock(500)
)

func increment(i int, m *ALock, wg *sync.WaitGroup) {
	for {
		index := m.Lock()
		if x >= 200 {
			m.UnLock(index)
			wg.Done()
			return
		}
		fmt.Println(x)
		x++
		m.UnLock(index)
	}
}

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
	for !a.flag[slot] {
	}
	return slot
}

func (a *ALock) UnLock(slot int) {
	a.flag[slot] = false
	a.flag[(slot+1)%a.size] = true
}
