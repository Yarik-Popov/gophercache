package cache

import (
	"flag"
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
	peers := strings.Split(peerAddresses, ",")

	config := Config{
		MaxElements:   maxElements,
		LocalAddress:  localAddress,
		ExpirySeconds: seconds,
		PeerAddresses: peers,
	}
	return &config, nil
}
