# gophercache
Distributed key value store written in Go

- Built a thread-safe, generic LRU cache in Go with configurable capacity and TTL-based entry expiration, leveraging Go generics (comparable/any type constraints) and a doubly-linked list + hash map for O(1) get/put operations
- Designed and implemented a distributed key-value store in Go featuring an LRU eviction policy, mutex-based concurrency control, and sliding TTL expiration to ensure data freshness under concurrent load

- Will expose cache operations via a RESTful HTTP API (GET /key, POST /key) using Go's standard net/http package, enabling language-agnostic client integration with zero external dependencies
- Will implement peer discovery and horizontal scalability using consistent hashing to distribute keys across cache nodes, minimizing key remapping during node additions or failures in the distributed cluster

