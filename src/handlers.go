package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func StartServer() {
	// Config
	config, err := CreateConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Setup
	server, err := CreateServer(config)
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
			log.Printf("server %s closed", config.LocalAddress)
		} else if err != nil {
			log.Printf("error listening on server %s: %s\n", config.LocalAddress, err)
		}
		cancelCtx()
	}()

	log.Println("Starting server")
	server.Print()
	<-ctx.Done()

}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	log.Println("Got /get/", key)

	ownerAddr, err := s.GetOwner(key)
	if err != nil {
		log.Print(err)
		http.Error(w, "Unexpected error", 500)
		return
	}

	if s.LocalAddress == ownerAddr {
		keyValueStore := s.localCache
		val, ok := keyValueStore.Get(key)
		writeGetResponse(w, key, val, ok)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/get/%s", ownerAddr, key))
	if err != nil {
		log.Print(err)
		http.Error(w, "Unexpected error", 500)
		return
	}

	defer resp.Body.Close()
	val, err := io.ReadAll(resp.Body)
	writeGetResponse(w, key, val, resp.StatusCode == http.StatusOK)
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	log.Println("Got /put/", key)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read body: %s\n", err)
		w.WriteHeader(400)
		return
	}

	ownerAddr, err := s.GetOwner(key)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unexpected error", 500)
		return
	}

	if s.LocalAddress == ownerAddr {
		keyValueStore := s.localCache
		keyValueStore.Put(key, body)

		log.Printf("Updated key '%s' with '%s'\n", key, body)
		w.WriteHeader(200)
		io.WriteString(w, "Succesfully updated key")
		return
	}

	log.Println("Redirecting to ", ownerAddr)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/put/%s", ownerAddr, key), bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		http.Error(w, "Unexpected error", 500)
		return
	}

	// Set the appropriate Content-Type header if needed
	req.Header.Set("Content-Type", "text/plain")

	// Use an http.Client to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unexpected error", 500)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		http.Error(w, req.Response.Status, req.Response.StatusCode)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	w.WriteHeader(200)
	w.Write(body)
	log.Println(string(body))
}

// Private functions

func writeGetResponse(w http.ResponseWriter, key string, val []byte, ok bool) {
	if ok {
		log.Printf("Found value: '%s' for key: '%s'\n", val, key)
		w.WriteHeader(200)
		w.Write(val)
	} else {
		errMsg := fmt.Sprintf("Key '%s' not found\n", key)
		log.Print(errMsg)
		http.Error(w, errMsg, 404)
	}
}
