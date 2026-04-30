package cache

import (
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	Cache *Cache[string, []byte]
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Got /get/", key)

	keyValueStore := s.Cache
	val, ok := keyValueStore.Get(key)
	if !ok {
		fmt.Printf("Key '%s' not found\n", key)
		w.WriteHeader(404)
		io.WriteString(w, "(nil)")
		return
	}

	fmt.Printf("Found value: '%s' for key: '%s'\n", val, key)
	w.WriteHeader(200)
	w.Write(val)
}

func (s *Server) HandlePut(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	fmt.Println("Got /put/", key)

	val, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Could not read body: %s\n", err)
		w.WriteHeader(400)
		return
	}

	keyValueStore := s.Cache
	keyValueStore.Put(key, val)

	fmt.Printf("Updated key '%s' with '%s'\n", key, val)
	w.WriteHeader(200)
	io.WriteString(w, "Succesfully updated key")
}
