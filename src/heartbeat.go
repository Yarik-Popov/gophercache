package cache

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type HeartBeater struct {
	MaxFailedHeartbeats uint
	NeighbourAddress    string
	RemainingAttempts   uint
	LocalAddress        string
	client              *http.Client
}

func CreateHeartBeater(
	address string,
	config *Config,
) *HeartBeater {
	client := &http.Client{
		Timeout: config.HeartbeatInterval,
	}

	hb := &HeartBeater{
		NeighbourAddress:    address,
		MaxFailedHeartbeats: config.MaxFailedHeartbeats,
		RemainingAttempts:   config.MaxFailedHeartbeats,
		LocalAddress:        config.LocalAddress,
		client:              client,
	}
	return hb
}

func (hb *HeartBeater) heartbeatPeer() bool {
	log.Printf("Sending heartbeat to %s", hb.NeighbourAddress)
	payload := []byte(hb.LocalAddress)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/heartbeat", hb.NeighbourAddress), bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	_, err = hb.client.Do(req)
	return err == nil
}

func (hb *HeartBeater) Heartbeat() bool {
	log.Printf("Attempting to heartbeat %s", hb.NeighbourAddress)
	if hb.RemainingAttempts <= 0 {
		log.Printf("Heart beat for %s failed", hb.NeighbourAddress)
		return false
	}

	res := hb.heartbeatPeer()
	if res {
		hb.RemainingAttempts = hb.MaxFailedHeartbeats
	} else {
		hb.RemainingAttempts--
	}

	log.Printf("%d remaining attempts left for %s", hb.RemainingAttempts, hb.NeighbourAddress)
	return true
}
