package main

import (
	"Yarik-Popov/gophercache/src"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// TODO: Make this configurable
	address := ":8080"
	var maxElements uint32 = 3
	ttl := 10 * time.Second

	keyValueStore := cache.CreateCache[string, []byte](maxElements, ttl)
	server := cache.Server{Cache: keyValueStore}
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /get/{key}", server.HandleGet)
	mux.HandleFunc("PUT /put/{key}", server.HandlePut)

	// Create server
	ctx, cancelCtx := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	// Start server
	go func() {
		err := httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server %s closed", address)
		} else if err != nil {
			fmt.Printf("error listening on server %s: %s\n", address, err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
