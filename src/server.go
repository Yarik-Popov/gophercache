package cache

import (
	"errors"
	"fmt"
	"hash/crc64"
	"math"
	"slices"
)

type Server struct {
	LocalAddress  string
	Peers         []string
	localCache    *Cache[string, []byte]
	ringOrdering  []uint64
	hashesToPeers map[uint64]string
	table         *crc64.Table // Needed to compute crc64
}

func CreateServer(config *Config) (*Server, error) {
	keyValueStore := CreateCache[string, []byte](config.MaxElements, config.ExpirySeconds)
	localAddress := config.LocalAddress
	peers := config.PeerAddresses
	numNodes := 1 + len(peers)

	server := Server{
		localCache:    keyValueStore,
		LocalAddress:  localAddress,
		Peers:         peers,
		ringOrdering:  make([]uint64, numNodes),
		hashesToPeers: make(map[uint64]string),
		table:         crc64.MakeTable(1000),
	}

	for i, peer := range peers {
		peerHash := crc64.Checksum([]byte(peer), server.table)
		server.ringOrdering[i] = peerHash
		server.hashesToPeers[peerHash] = peer
		fmt.Printf("i: %d, peer:%s, peerHash: %d\n", i, peer, peerHash)
	}

	peerHash := crc64.Checksum([]byte(localAddress), server.table)
	server.ringOrdering[numNodes-1] = peerHash
	server.hashesToPeers[peerHash] = localAddress
	fmt.Printf("numNodes-1: %d, peer:%s, peerHash: %d\n", numNodes-1, localAddress, peerHash)

	slices.Sort(server.ringOrdering)
	return &server, nil
}

func (s *Server) Print() {
	fmt.Printf("Server on %s", s.LocalAddress)
	ttl := s.localCache.duration
	if ttl == 0 {
		fmt.Print(" with no ttl")
	} else {
		fmt.Printf(" with ttl of %d", ttl)
	}

	fmt.Println(" with ordering:")

	// First node is responsible for everything after the last node and before the first
	var prevHash uint64 = 0
	for _, hash := range s.ringOrdering {
		nodeName := s.hashesToPeers[hash]
		fmt.Printf("\t%s is responsible for [%d, %d]\n", nodeName, prevHash, hash)
		prevHash = hash + 1
	}

	// Print wrap around node
	hash := s.ringOrdering[0]
	nodeName := s.hashesToPeers[hash]
	fmt.Printf("\t%s is responsible for [%d, %d]\n", nodeName, prevHash, uint64(math.MaxUint64))

}

func (s *Server) GetOwner(key string) (string, error) {
	keyHash := crc64.Checksum([]byte(key), s.table)
	nodeIndex, _ := slices.BinarySearch(s.ringOrdering, keyHash)

	var nodeHash uint64
	if nodeIndex >= len(s.ringOrdering) {
		nodeHash = s.ringOrdering[0]
	} else {
		nodeHash = s.ringOrdering[nodeIndex]
	}

	nodeAddress, ok := s.hashesToPeers[nodeHash]
	if !ok {
		return "", errors.New("Node not found")
	}
	return nodeAddress, nil
}

// TODO:(8)
func (s *Server) IsOwner(key string) bool {
	owner, err := s.GetOwner(key)
	return err == nil && owner == s.LocalAddress
}
