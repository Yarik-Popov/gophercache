package src

import (
	"Yarik-Popov/gophercache/src"
	"testing"
	"time"
)

func TestCache_TTLEviction(t *testing.T) {
	// Create cache with 100ms TTL
	duration := 100 * time.Millisecond
	c := cache.CreateCache[string, string](10, duration)

	c.Put("ephemeral", "I vanish")

	// 1. Immediate access should work
	if _, ok := c.Get("ephemeral"); !ok {
		t.Error("Expected key to be present immediately after Put")
	}

	// 2. Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)

	// 3. Access after expiration should fail and trigger discard
	if _, ok := c.Get("ephemeral"); ok {
		t.Error("Expected key to be discarded after TTL expiration")
	}

	// 4. Verify internal state was cleaned up
	if c.NumElements != 0 {
		t.Errorf("Expected numElements to be 0 after lazy eviction, got %d", c.NumElements)
	}
}

func TestCache_TTLDisabled(t *testing.T) {
	// Duration 0 should disable TTL
	c := cache.CreateCache[string, string](10, 0)

	c.Put("immortal", "I stay")

	// Wait a bit to simulate passage of time
	time.Sleep(50 * time.Millisecond)

	if val, ok := c.Get("immortal"); !ok || val != "I stay" {
		t.Error("Expected key to persist when TTL is disabled (duration=0)")
	}
}

func TestCache_TTLResetOnPut(t *testing.T) {
	duration := 100 * time.Millisecond
	c := cache.CreateCache[string, string](10, duration)

	c.Put("refresh", "v1")

	// Wait halfway through TTL
	time.Sleep(60 * time.Millisecond)

	// Overwrite key - this should reset the deadline to +100ms from now
	c.Put("refresh", "v2")

	// Wait 60ms more.
	// Total time since first Put: 120ms (would have expired)
	// Total time since second Put: 60ms (should still be alive)
	time.Sleep(60 * time.Millisecond)

	if val, ok := c.Get("refresh"); !ok || val != "v2" {
		t.Error("Expected Put to reset the TTL deadline")
	}
}
