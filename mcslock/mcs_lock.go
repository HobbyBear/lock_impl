package main

import (
	"fmt"
	"lock_impl"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type MCSLock struct {
	tail unsafe.Pointer
}

func NewMCSLock() *MCSLock {
	n := unsafe.Pointer(nil)
	return &MCSLock{tail: n}
}

func (c *MCSLock) Lock(myNode *lock_impl.QNode) {
	myNode.Locked = true
	preNode := (*lock_impl.QNode)(atomic.SwapPointer(&c.tail, unsafe.Pointer(myNode)))
	if preNode != nil {
		myNode.Locked = true
		preNode.Next = myNode
		for myNode.Locked {
			runtime.Gosched()
		}
	}

}

func (c *MCSLock) UnLock(myNode *lock_impl.QNode) {
	myNode.Locked = false
	if myNode.Next == nil {
		if atomic.CompareAndSwapPointer(&c.tail, unsafe.Pointer(myNode), nil) {
			return
		}
		for myNode.Next == nil {
		}
	}
	myNode.Next.Locked = false
	myNode.Next = nil
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
	m = NewMCSLock()
)

func increment(i int, m *MCSLock, wg *sync.WaitGroup) {
	myNode := new(lock_impl.QNode)

	for {
		m.Lock(myNode)
		if x >= 200 {
			m.UnLock(myNode)
			wg.Done()
			return
		}
		fmt.Println(x)
		x++
		m.UnLock(myNode)
	}
}
