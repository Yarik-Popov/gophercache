package cache

import (
	"flag"
	"strings"
	"time"
)

type Config struct {
	MaxElements          uint
	LocalAddress         string
	ExpirySeconds        time.Duration
	PeerAddresses        []string
	InitialHeartbeatWait time.Duration
	HeartbeatInterval    time.Duration
	MaxFailedHeartbeats  uint
}

func CreateConfig() (*Config, error) {
	maxElementsPtr := flag.Uint("maxelements", 5, "Max number of elements to be stored.")
	localAddressPtr := flag.String("address", "http://localhost:8080", "Address to run server on.")
	expirySecondsPtr := flag.Uint("ttl", 10, "Expiry time in seconds. 0 disables expiry.")
	peerAddressesPtr := flag.String("peers", "", "Comma separated list of peer addresses")

	// Heartbeat
	initialHeartbeatWaitPtr := flag.Uint("init", 10, "Initial number of seconds before starting to heartbeat peers")
	heartbeatIntervalPtr := flag.Uint("interval", 10, "Interval in seconds to heartbeat peers")
	maxFailedHeartbeatsPtr := flag.Uint("heartbeatfailures", 3, "Maximum number of heartbeats that can be missed before declaring node to be dead")

	flag.Parse()

	maxElements := *maxElementsPtr
	localAddress := *localAddressPtr
	peerAddresses := *peerAddressesPtr

	expirySeconds := *expirySecondsPtr
	seconds := time.Duration(expirySeconds) * time.Second

	initialHeartbeatWait := *initialHeartbeatWaitPtr
	initialHeartbeatWaitSeconds := time.Duration(initialHeartbeatWait) * time.Second

	heartbeatInterval := *heartbeatIntervalPtr
	heartbeatIntervalSeconds := time.Duration(heartbeatInterval) * time.Second

	maxFailedHeartbeats := *maxFailedHeartbeatsPtr

	// string.Split returns an array with the first element being the input if it can't split
	var peers []string
	if peerAddresses != "" {
		peers = strings.Split(peerAddresses, ",")
	}

	config := Config{
		MaxElements:          maxElements,
		LocalAddress:         localAddress,
		ExpirySeconds:        seconds,
		PeerAddresses:        peers,
		InitialHeartbeatWait: initialHeartbeatWaitSeconds,
		HeartbeatInterval:    heartbeatIntervalSeconds,
		MaxFailedHeartbeats:  maxFailedHeartbeats,
	}
	return &config, nil
}
