package cache

import (
	"time"
)

type nodeElement[K comparable, V any] struct {
	Prev       *nodeElement[K, V]
	Next       *nodeElement[K, V]
	Value      V
	Key        K
	ExpiryTime time.Time
}

func (node *nodeElement[K, V]) removeNode() {
	prev := node.Prev
	next := node.Next
	prev.Next = next
	next.Prev = prev
}
