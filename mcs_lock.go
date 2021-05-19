package lock_impl

import (
	"sync/atomic"
	"unsafe"
)

type MCSLock struct {
	tail unsafe.Pointer
}

func NewMCSLock() *MCSLock {
	 n := unsafe.Pointer(nil)
	return &MCSLock{tail:n}
}

func (c *MCSLock)Lock(myNode *QNode)  {
	myNode.Locked = true
	preNode := (*QNode)(atomic.SwapPointer(&c.tail, unsafe.Pointer(myNode)))
	if  preNode != nil {
		myNode.Locked = true
		preNode.Next = myNode
		for myNode.Locked{}
	}

}

func (c *MCSLock) UnLock(myNode *QNode) {
	myNode.Locked =false
	if myNode.Next == nil{
		if atomic.CompareAndSwapPointer(&c.tail,unsafe.Pointer(myNode),nil){
			return
		}
		for myNode.Next == nil{}
	}
	myNode.Next.Locked = false
	myNode.Next = nil
}