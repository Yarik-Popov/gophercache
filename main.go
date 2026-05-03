package main

import (
	"Yarik-Popov/gophercache/src"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	// httpServer doesn't like that we start with the protocol
	addr := config.LocalAddress
	if strings.HasPrefix(config.LocalAddress, "http://") {
		addr = config.LocalAddress[len("http://"):]
	} else if strings.HasPrefix(config.LocalAddress, "https://") {
		addr = config.LocalAddress[len("https://"):]
	}

	httpServer := &http.Server{
		Addr:    addr,
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

	fmt.Println("Starting server")
	server.Print()
	<-ctx.Done()
}
