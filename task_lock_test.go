package lock_impl

import "testing"

func TestTasLock_Lock(t *testing.T) {
	l := TTasLock{State: new(int32)}
	l.Lock()
	l.UnLock()
}
