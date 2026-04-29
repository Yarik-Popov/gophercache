package src

import (
	"Yarik-Popov/gophercache/src"
	"sync"
	"testing"
)

// TestCache_ConcurrentPuts verifies that concurrent Put operations do not cause
// data races or corrupt internal state.
func TestCache_ConcurrentPuts(t *testing.T) {
	c := cache.CreateCache[int, int](100, 0)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(n, n*10)
		}(i)
	}
	wg.Wait()

	if c.NumElements > c.MaxElements {
		t.Errorf("NumElements (%d) exceeds MaxElements (%d) after concurrent puts", c.NumElements, c.MaxElements)
	}
}

// TestCache_ConcurrentGets verifies that concurrent Get operations do not cause
// data races when reading the same keys.
func TestCache_ConcurrentGets(t *testing.T) {
	c := cache.CreateCache[string, int](10, 0)
	c.Put("shared", 42)

	type result struct {
		val int
		ok  bool
	}
	results := make(chan result, 50)

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, ok := c.Get("shared")
			results <- result{val, ok}
		}()
	}
	wg.Wait()
	close(results)

	for r := range results {
		if !r.ok || r.val != 42 {
			t.Errorf("expected (42, true), got (%v, %v)", r.val, r.ok)
		}
	}
}

// TestCache_ConcurrentMixedOperations verifies that interleaved Puts and Gets
// from multiple goroutines do not cause data races or panics.
func TestCache_ConcurrentMixedOperations(t *testing.T) {
	c := cache.CreateCache[int, int](50, 0)

	var wg sync.WaitGroup

	// Writers
	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(n, n)
		}(i)
	}

	// Readers (keys may or may not exist yet — that is fine)
	for i := range 50 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Get(n)
		}(i)
	}

	wg.Wait()

	if c.NumElements > c.MaxElements {
		t.Errorf("NumElements (%d) exceeds MaxElements (%d) after mixed concurrent ops", c.NumElements, c.MaxElements)
	}
}

// TestCache_ConcurrentEviction verifies that the LRU eviction path is safe
// under concurrent pressure: writes that exceed capacity must never leave the
// cache in an inconsistent state.
func TestCache_ConcurrentEviction(t *testing.T) {
	const MaxElems = 10
	c := cache.CreateCache[int, int](MaxElems, 0)

	var wg sync.WaitGroup
	// Write more keys than the capacity to force eviction on every goroutine.
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(n, n)
		}(i)
	}
	wg.Wait()

	if c.NumElements > c.MaxElements {
		t.Errorf("NumElements (%d) exceeds MaxElements (%d) after concurrent eviction", c.NumElements, c.MaxElements)
	}

	// Sentinel nodes must still be consistent.
	if c.First() == nil || c.Last() == nil {
		t.Error("sentinel node pointers are nil after concurrent eviction")
	}
}

// TestCache_ConcurrentUpdates verifies that updating the same key from multiple
// goroutines simultaneously does not cause a data race or panic.
func TestCache_ConcurrentUpdates(t *testing.T) {
	c := cache.CreateCache[string, int](5, 0)
	c.Put("counter", 0)

	var wg sync.WaitGroup
	for i := 1; i <= 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put("counter", n)
		}(i)
	}
	wg.Wait()

	// The key must still be present; any of the written values is acceptable.
	if _, ok := c.Get("counter"); !ok {
		t.Error("expected 'counter' key to be present after concurrent updates")
	}
}
