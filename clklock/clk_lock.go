package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type CLKLock struct {
	tail unsafe.Pointer
}

type QNode struct {
	Locked bool
}

func NewCLKLock() *CLKLock {
	node := unsafe.Pointer(new(QNode))
	return &CLKLock{
		tail: node,
	}
}

type container struct {
	MyNode  unsafe.Pointer
	PreNode unsafe.Pointer
}

func (c *CLKLock) Lock(contain *container) {
	(*QNode)(contain.MyNode).Locked = true
	contain.PreNode = atomic.SwapPointer(&c.tail, contain.MyNode)
	//fmt.Println("tail", *(*QNode)(c.tail))
	for (*QNode)(contain.PreNode).Locked {
		runtime.Gosched()
	}
}

func (c *CLKLock) UnLock(contain *container) {
	(*QNode)(contain.MyNode).Locked = false
	//fmt.Println("结束", "tail", *(*QNode)(c.tail))
	contain.MyNode = contain.PreNode
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
	m = NewCLKLock()
)

func increment(i int, m *CLKLock, wg *sync.WaitGroup) {
	contain := &container{
		MyNode:  unsafe.Pointer(new(QNode)),
		PreNode: unsafe.Pointer(nil),
	}
	for {
		m.Lock(contain)
		if x >= 200 {
			m.UnLock(contain)
			wg.Done()
			return
		}
		fmt.Println(x)
		x++
		m.UnLock(contain)
	}
}
