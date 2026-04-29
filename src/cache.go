package cache

import (
	"sync"
	"time"
)

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
}

type Cache[K comparable, V any] struct {
	MaxElements uint32
	NumElements uint32
	front       *nodeElement[K, V]
	back        *nodeElement[K, V]
	lookup      map[K]*nodeElement[K, V]
	duration    time.Duration
	lock        sync.Mutex
}

func CreateCache[K comparable, V any](maxElements uint32, duration time.Duration) *Cache[K, V] {
	// Setup sentinel front and back nodes to make life easier when moving elements around
	front := &nodeElement[K, V]{
		Prev: nil,
		Next: nil,
	}
	back := &nodeElement[K, V]{
		Prev: front,
		Next: nil,
	}
	front.Next = back

	cache := &Cache[K, V]{
		MaxElements: maxElements,
		NumElements: 0,
		front:       front,
		back:        back,
		lookup:      make(map[K]*nodeElement[K, V]),
		duration:    duration,
	}
	return cache
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	node, ok := c.lookup[key]
	if !ok {
		var zero V // Gets zeroed out so still has a value. Now code compiles
		return zero, false
	}

	// This case should never happen but just in case panic
	if node == nil {
		panic("Node is expected but not found")
	}

	if c.duration > 0 {
		now := time.Now()
		if node.ExpiryTime.Before(now) {
			c.deleteNode(node)
			return node.Value, false
		}
		node.ExpiryTime = now.Add(c.duration)
	}

	node.removeNode()
	c.insertNode(node)

	return node.Value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()

	node, ok := c.lookup[key]
	if ok {
		node.Value = value
		node.ExpiryTime = time.Now().Add(c.duration)

		node.removeNode()
		c.insertNode(node)
		return
	}

	newNode := &nodeElement[K, V]{
		Value:      value,
		Key:        key,
		ExpiryTime: time.Now().Add(c.duration),
	}

	if c.MaxElements == c.NumElements {
		oldLast := c.back.Prev
		c.deleteNode(oldLast)
	}

	c.NumElements++
	c.insertNode(newNode)
	c.lookup[key] = newNode
}

func (c *Cache[K, V]) First() *nodeElement[K, V] {
	return c.front.Next
}

func (c *Cache[K, V]) Last() *nodeElement[K, V] {
	return c.back.Prev
}

// Private

func (c *Cache[K, V]) insertNode(node *nodeElement[K, V]) {
	oldFirst := c.front.Next
	c.front.Next = node
	node.Prev = c.front
	node.Next = oldFirst
	oldFirst.Prev = node
}

func (c *Cache[K, V]) deleteNode(node *nodeElement[K, V]) {
	node.removeNode()
	c.NumElements--
	delete(c.lookup, node.Key)
}
