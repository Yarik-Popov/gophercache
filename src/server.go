package cache

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"slices"
)

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

type Server struct {
	config        *Config
	LocalAddress  string
	Peers         []string
	localCache    *Cache[string, []byte]
	ringOrdering  []uint64
	hashesToPeers map[uint64]string
	currentIndex  int
}

func CreateServer(config *Config) (*Server, error) {
	keyValueStore := CreateCache[string, []byte](config.MaxElements, config.ExpirySeconds)
	localAddress := config.LocalAddress
	peers := config.PeerAddresses
	numNodes := 1 + len(peers)

	server := Server{
		config:        config,
		localCache:    keyValueStore,
		LocalAddress:  localAddress,
		Peers:         peers,
		ringOrdering:  make([]uint64, numNodes),
		hashesToPeers: make(map[uint64]string),
	}

	for i, peer := range peers {
		peerHash := hashString(peer)
		server.ringOrdering[i] = peerHash
		server.hashesToPeers[peerHash] = peer
	}

	currentNodeHash := hashString(localAddress)
	server.ringOrdering[numNodes-1] = currentNodeHash
	server.hashesToPeers[currentNodeHash] = localAddress

	slices.Sort(server.ringOrdering)

	idx := slices.Index(server.ringOrdering, currentNodeHash)
	if idx == -1 {
		return nil, errors.New("Cannot find current node in ring ordering")
	}
	server.currentIndex = idx

	return &server, nil
}

func (s *Server) PrevServer() string {
	idx := s.currentIndex - 1
	if idx < 0 {
		idx = len(s.ringOrdering) - 1
	}
	return s.hashesToPeers[s.ringOrdering[idx]]
}

func (s *Server) NextServer() string {
	idx := s.currentIndex + 1
	if idx == len(s.ringOrdering) {
		idx = 0
	}
	return s.hashesToPeers[s.ringOrdering[idx]]
}

func (s *Server) Print() {
	fmt.Printf("Server on %s", s.LocalAddress)
	ttl := s.localCache.duration
	if ttl == 0 {
		fmt.Print(" with unlimited ttl")
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
	keyHash := hashString(key)
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
