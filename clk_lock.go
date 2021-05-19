package lock_impl

import (
	"sync/atomic"
	"unsafe"
)

type CLKLock struct {
	tail unsafe.Pointer
}

func NewCLKLock() *CLKLock {
	node := unsafe.Pointer(new(QNode))
	return &CLKLock{
		tail: node,
	}
}

func (c *CLKLock)Lock(myNode *QNode,preNode *QNode)  {
	myNode.Locked = true
	preNode = (*QNode)(atomic.SwapPointer(&c.tail, unsafe.Pointer(myNode)))
	for preNode.Locked {}
}

func (c *CLKLock) UnLock(myNode *QNode,preNode *QNode) {
	myNode.Locked =false
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(preNode)),unsafe.Pointer(myNode))
}