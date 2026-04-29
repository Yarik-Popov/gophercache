package src

import (
	"Yarik-Popov/gophercache/src"
	"testing"
)

func TestCache_BasicOperations(t *testing.T) {
	// Initialize with capacity for 2
	c := cache.CreateCache[string, int](2, 0)

	// Test 1: Get a non-existent key
	if _, ok := c.Get("missing"); ok {
		t.Errorf("expected ok=false for missing key")
	}

	// Test 2: Put and Get
	c.Put("a", 1)
	val, ok := c.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	// Test 3: Update existing key
	c.Put("a", 10)
	val, ok = c.Get("a")
	if !ok || val != 10 {
		t.Errorf("expected 10 after update, got %v", val)
	}
}

func TestCache_EvictionLogic(t *testing.T) {
	// Capacity of 2
	c := cache.CreateCache[string, string](2, 0)

	c.Put("first", "A")  // [A]
	c.Put("second", "B") // [B, A]

	// Access "first" to make it Most Recently Used (MRU)
	c.Get("first") // [A, B]

	// This should trigger eviction of "second" because "first" was recently accessed
	c.Put("third", "C") // [C, A]

	if _, ok := c.Get("second"); ok {
		t.Error("expected 'second' to be evicted")
	}

	if val, ok := c.Get("first"); !ok || val != "A" {
		t.Errorf("expected 'first' to remain in cache, got %v", val)
	}

	if c.NumElements > c.MaxElements {
		t.Errorf("numElements (%d) exceeds maxElements (%d)", c.NumElements, c.MaxElements)
	}
}

func TestCache_Generics(t *testing.T) {
	t.Run("IntKeys_StructValues", func(t *testing.T) {
		type User struct{ Name string }
		c := cache.CreateCache[int, User](1, 0)

		user := User{Name: "Alice"}
		c.Put(1, user)

		got, _ := c.Get(1)
		if got.Name != "Alice" {
			t.Errorf("Generics failed: expected Alice, got %s", got.Name)
		}
	})
}

func TestCache_InternalPointers(t *testing.T) {
	// This test ensures your doubly linked list isn't breaking
	c := cache.CreateCache[string, int](3, 0)

	c.Put("1", 1)
	c.Put("2", 2)
	c.Put("3", 3)

	// Current order (Front to Back): 3 -> 2 -> 1
	if c.First().Key != "3" || c.Last().Key != "1" {
		t.Errorf("Initial pointer state wrong: front=%v, back=%v", c.First().Key, c.Last().Key)
	}

	// Move middle to front
	c.Get("2") // Order: 2 -> 3 -> 1

	if c.First().Key != "2" {
		t.Errorf("Get didn't move key to front. Got front=%v", c.First().Key)
	}

	if c.First().Next.Key != "3" || c.Last().Prev.Key != "3" {
		t.Error("Middle pointers broken after Get")
	}
}
