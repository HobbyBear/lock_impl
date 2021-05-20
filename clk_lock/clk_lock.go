package main

import (
	"fmt"
	"lock_impl"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type CLKLock struct {
	tail unsafe.Pointer
}

func NewCLKLock() *CLKLock {
	node := unsafe.Pointer(new(lock_impl.QNode))
	return &CLKLock{
		tail: node,
	}
}

type container struct {
	MyNode  unsafe.Pointer
	PreNode unsafe.Pointer
}

func (c *CLKLock) Lock(contain *container, num int) {
	(*lock_impl.QNode)(atomic.LoadPointer(&contain.MyNode)).Locked = true
	//(*QNode)(contain.MyNode).Locked = true
	//fmt.Println(contain.MyNode,num,"lock")
	contain.PreNode = atomic.SwapPointer(&c.tail, contain.MyNode)
	//fmt.Println(unsafe.Pointer(contain.PreNode),num)
	//fmt.Println(c.tail)
	//fmt.Println(unsafe.Pointer(preNode))
	//i := 0
	for (*lock_impl.QNode)(atomic.LoadPointer(&contain.PreNode)).Locked {
		runtime.Gosched()
	}
	//fmt.Println(num)
}

func (c *CLKLock) UnLock(contain *container, num int) {
	//(*QNode)(contain.MyNode).Locked = false
	(*lock_impl.QNode)(atomic.LoadPointer(&contain.MyNode)).Locked = false
	//fmt.Println("my", unsafe.Pointer(contain.MyNode))
	//fmt.Println((*QNode)(atomic.LoadPointer(&c.tail)))
	//atomic.LoadPointer(&c.tail)
	//fmt.Println(unsafe.Pointer(contain.MyNode),num,"unlock")
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
	myNode := new(lock_impl.QNode)
	//old := new(QNode)
	c := new(container)
	c.MyNode = unsafe.Pointer(myNode)
	c.PreNode = unsafe.Pointer(nil)
	for {
		m.Lock(c, i)
		if x >= 200 {
			m.UnLock(c, i)
			wg.Done()
			return
		}
		fmt.Println(x)
		x++
		m.UnLock(c, i)
	}
}
