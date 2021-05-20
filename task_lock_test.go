package lock_impl

import (
	"testing"
)

func TestTasLock_Lock(t *testing.T) {

	l := TTasLock{
		State: new(int32),
	}
	num := 0
	for i := 0; i < 1000; i++ {
		go func() {
			l.Lock()
			num++
			l.UnLock()
		}()
	}

}
