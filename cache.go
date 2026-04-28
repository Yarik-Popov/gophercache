package main

import (
	"sync"
	"time"
)

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
}

type Cache[K comparable, V any] struct {
	maxElements uint32
	numElements uint32
	front       *nodeElement[K, V]
	back        *nodeElement[K, V]
	lookup      map[K]*nodeElement[K, V]
	duration    time.Duration
	lock        sync.Mutex
}

func CreateCache[K comparable, V any](maxElements uint32, duration time.Duration) *Cache[K, V] {
	// Setup sentinel front and back nodes to make life easier when moving elements around
	front := &nodeElement[K, V]{
		prev: nil,
		next: nil,
	}
	back := &nodeElement[K, V]{
		prev: front,
		next: nil,
	}
	front.next = back

	cache := &Cache[K, V]{
		maxElements: maxElements,
		numElements: 0,
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
		if node.expiryTime.Before(now) {
			c.deleteNode(node)
			return node.value, false
		}
		node.expiryTime = now.Add(c.duration)
	}

	node.removeNode()
	c.insertNode(node)

	return node.value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	c.lock.Lock()
	defer c.lock.Unlock()

	node, ok := c.lookup[key]
	if ok {
		node.value = value
		node.expiryTime = time.Now().Add(c.duration)

		node.removeNode()
		c.insertNode(node)
		return
	}

	newNode := &nodeElement[K, V]{
		value:      value,
		key:        key,
		expiryTime: time.Now().Add(c.duration),
	}

	if c.maxElements == c.numElements {
		oldLast := c.back.prev
		c.deleteNode(oldLast)
	}

	c.numElements++
	c.insertNode(newNode)
	c.lookup[key] = newNode
}

// Private

type nodeElement[K comparable, V any] struct {
	prev       *nodeElement[K, V]
	next       *nodeElement[K, V]
	value      V
	key        K
	expiryTime time.Time
}

func (c *Cache[K, V]) insertNode(node *nodeElement[K, V]) {
	oldFirst := c.front.next
	c.front.next = node
	node.prev = c.front
	node.next = oldFirst
	oldFirst.prev = node
}

func (c *Cache[K, V]) deleteNode(node *nodeElement[K, V]) {
	node.removeNode()
	c.numElements--
	delete(c.lookup, node.key)
}

func (node *nodeElement[K, V]) removeNode() {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}
