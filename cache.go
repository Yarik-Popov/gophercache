package main

type Cache[K comparable, V any] struct {
	maxElements uint32
	numElements uint32
	front       *nodeElement[K, V]
	back        *nodeElement[K, V]
	lookup      map[K]*nodeElement[K, V]
}

func CreateCache[K comparable, V any](maxElements uint32) *Cache[K, V] {
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
	}
	return cache
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	node, ok := c.lookup[key]
	if !ok {
		return node.value, false
	}

	// This case should never happen but just in case panic
	if node == nil {
		panic("Node is expected but not found")
	}

	node.removeNode()
	c.insertNode(node)

	return node.value, true
}

func (c *Cache[K, V]) Put(key K, value V) {
	node, ok := c.lookup[key]
	if ok {
		node.value = value

		node.removeNode()
		c.insertNode(node)
		return
	}

	if c.maxElements == c.numElements {
		oldLast := c.back.prev
		oldLast.removeNode()
		c.numElements--
		delete(c.lookup, oldLast.key)
	}

	c.numElements++
	newNode := &nodeElement[K, V]{
		value: value,
		key:   key,
	}
	c.insertNode(newNode)
	c.lookup[key] = newNode
}

// Private

type nodeElement[K comparable, V any] struct {
	prev  *nodeElement[K, V]
	next  *nodeElement[K, V]
	value V
	key   K
}

func (c *Cache[K, V]) insertNode(node *nodeElement[K, V]) {
	oldFirst := c.front.next
	c.front.next = node
	node.prev = c.front
	node.next = oldFirst
	oldFirst.prev = node
}

func (node *nodeElement[K, V]) removeNode() {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}
