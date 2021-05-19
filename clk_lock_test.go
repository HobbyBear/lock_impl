package lock_impl

import (
	"log"
	"testing"
	"time"
)

func TestCLKLock_Lock(t *testing.T) {
	l := NewMCSLock()

	go func() {
		myNode := new(QNode)
		l.Lock(myNode)
		time.Sleep(2 * time.Second)
		log.Println("go1")
		l.UnLock(myNode)
	}()

	go func() {
		myNode := new(QNode)
		l.Lock(myNode)
		log.Println("go2")
		l.UnLock(myNode)
	}()

	time.Sleep(5 * time.Second)
}
