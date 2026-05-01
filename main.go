package main

import (
	"Yarik-Popov/gophercache/src"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Config
	config, err := cache.CreateConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Setup
	server, err := cache.CreateServer(config)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /get/{key}", server.HandleGet)
	mux.HandleFunc("PUT /put/{key}", server.HandlePut)

	// Create server
	ctx, cancelCtx := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:    config.LocalAddress,
		Handler: mux,
	}

	// Start server
	go func() {
		err := httpServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server %s closed", config.LocalAddress)
		} else if err != nil {
			fmt.Printf("error listening on server %s: %s\n", config.LocalAddress, err)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
