package lock_impl

type QNode struct {
	Locked bool
	Next   *QNode
}
