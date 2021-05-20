package lock_impl

import (
	"fmt"
	"lock_impl/ttaslock"
	"sync"
	"testing"
)

func TestTasLock_Lock(t *testing.T) {

	l := ttaslock.TTasLock{
		State: new(int32),
	}
	num := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Lock()
			num++
			l.UnLock()
		}()
	}
	wg.Wait()
	fmt.Println(num)

}
