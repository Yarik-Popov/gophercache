package cache

import (
	"errors"
	"hash/crc64"
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
		table:         crc64.MakeTable(0),
	}

	for i, peer := range peers {
		peerHash := crc64.Checksum([]byte(peer), server.table)
		server.ringOrdering[i] = peerHash
		server.hashesToPeers[peerHash] = peer
	}

	peerHash := crc64.Checksum([]byte(localAddress), server.table)
	server.ringOrdering[numNodes-1] = peerHash
	server.hashesToPeers[peerHash] = localAddress

	slices.Sort(server.ringOrdering)
	return &server, nil
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

func (s *Server) IsOwner(key string) bool {
	owner, err := s.GetOwner(key)
	return err == nil && owner == s.LocalAddress
}
