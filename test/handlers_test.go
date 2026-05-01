package src

import (
	"Yarik-Popov/gophercache/src"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlers(t *testing.T) {
	// Setup: Create a keyValueStore and a server
	// Using a small TTL for testing
	config := &cache.Config{
		ExpirySeconds: 500 * time.Millisecond,
		MaxElements:   10,
	}
	server, _ := cache.CreateServer(config)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /get/{key}", server.HandleGet)
	mux.HandleFunc("PUT /put/{key}", server.HandlePut)

	t.Run("PUT /put/apple", func(t *testing.T) {
		body := `{"value": "red"}`
		req := httptest.NewRequest(http.MethodPut, "/put/apple", strings.NewReader(body))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusCreated {
			t.Errorf("expected 200/201, got %d", w.Code)
		}
	})

	t.Run("GET /get/apple (Found)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get/apple", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "red") {
			t.Errorf("expected body to contain 'red', got %s", w.Body.String())
		}
	})

	t.Run("GET /get/banana (Not Found)", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/get/banana", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestAPITTL(t *testing.T) {
	config := &cache.Config{
		ExpirySeconds: 100 * time.Millisecond,
		MaxElements:   10,
	}
	server, _ := cache.CreateServer(config)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /get/{key}", server.HandleGet)
	mux.HandleFunc("PUT /put/{key}", server.HandlePut)

	mux.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("PUT", "/put/expire", strings.NewReader(`{"v":"gone"}`)))

	time.Sleep(150 * time.Millisecond)

	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "/get/expire", nil))

	if w.Code != http.StatusNotFound {
		t.Errorf("API should return 404 for expired TTL, got %d", w.Code)
	}
}
