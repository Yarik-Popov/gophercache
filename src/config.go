package cache

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	MaxElements   uint
	LocalAddress  string
	ExpirySeconds time.Duration
	PeerAddresses []string
}

func CreateConfig() (*Config, error) {
	maxElementsPtr := flag.Uint("maxelements", 5, "Max number of elements to be stored.")
	localAddressPtr := flag.String("address", "localhost:8080", "Address to run server on.")
	expirySecondsPtr := flag.Uint("ttl", 10, "Expiry time in seconds. 0 disables expiry.")
	peerAddressesPtr := flag.String("peers", "", "Comma separated list of peer addresses")

	flag.Parse()

	maxElements := *maxElementsPtr
	localAddress := *localAddressPtr
	expirySeconds := *expirySecondsPtr
	peerAddresses := *peerAddressesPtr

	seconds := time.Duration(expirySeconds) * time.Second

	// string.Split returns an array with the first element being the input if it can't split
	var peers []string
	if peerAddresses != "" {
		peers = strings.Split(peerAddresses, ",")
	}
	fmt.Printf("peers: %d\n", len(peers))
	fmt.Printf("peers: %v\n", (peers))
	fmt.Printf("peerAddresses: %s\n", peerAddresses)

	config := Config{
		MaxElements:   maxElements,
		LocalAddress:  localAddress,
		ExpirySeconds: seconds,
		PeerAddresses: peers,
	}
	return &config, nil
}
