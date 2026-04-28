package main

type Cache struct {
	maxElements int
	numElements int
	front       *nodeElement
	back        *nodeElement
	lookup      map[string]*nodeElement
}

func CreateCache(maxElements int) *Cache {
	// Setup sentinel front and back nodes to make life easier when moving elements around
	front := &nodeElement{
		prev:  nil,
		next:  nil,
		value: nil,
	}
	back := &nodeElement{
		prev:  front,
		next:  nil,
		value: nil,
	}
	front.next = back

	cache := &Cache{
		maxElements: maxElements,
		numElements: 0,
		front:       front,
		back:        back,
		lookup:      make(map[string]*nodeElement),
	}
	return cache
}

func (c *Cache) Get(key string) (any, bool) {
	node, ok := c.lookup[key]
	if !ok {
		return nil, false
	}

	// This case should never happen but just in case panic
	if node == nil {
		panic("Node is expected but not found")
	}

	removeNode(node)
	c.insertNode(node)

	return node.value, true
}

func (c *Cache) Put(key string, value any) {
	node, ok := c.lookup[key]
	if ok {
		node.value = value

		removeNode(node)
		c.insertNode(node)
		return
	}

	if c.maxElements == c.numElements {
		oldLast := c.back.prev
		removeNode(oldLast)
		c.numElements--
		delete(c.lookup, oldLast.key)
	}

	c.numElements++
	newNode := &nodeElement{
		value: value,
		key:   key,
	}
	c.insertNode(newNode)
	c.lookup[key] = newNode
}

// Private

type nodeElement struct {
	prev  *nodeElement
	next  *nodeElement
	value any
	key   string
}

func (c *Cache) insertNode(node *nodeElement) {
	oldFirst := c.front.next
	c.front.next = node
	node.prev = c.front
	node.next = oldFirst
	oldFirst.prev = node
}

func removeNode(node *nodeElement) {
	prev := node.prev
	next := node.next
	prev.next = next
	next.prev = prev
}
