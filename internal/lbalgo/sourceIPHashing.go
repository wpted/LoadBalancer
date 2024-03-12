package lbalgo

import (
    "LoadBalancer/internal/model"
    "hash/fnv"
    "math/rand"
    "net/http"
    "sync"
    "time"
)

// SIH is the struct used for source IP hashing.
type SIH struct {
    bucket map[int]map[string]struct{}
    rand   *rand.Rand
    sync.RWMutex
}

// NewSIH creates a SIH instance.
func NewSIH(backendServers *model.BEServers) *SIH {
    sih := &SIH{
        bucket: make(map[int]map[string]struct{}),
        rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
    }

    for i := 0; i < 10; i++ {
        sih.bucket[i] = make(map[string]struct{})
    }

    if backendServers != nil {
        for addr := range *backendServers {
            bucketNum := ihash(addr) % len(sih.bucket)
            sih.bucket[bucketNum][addr] = struct{}{}
        }
    }

    return sih
}

// ChooseServer chooses a server based on the clientIP.
func (s *SIH) ChooseServer(req *http.Request) (string, error) {
    clientIP := getClientIP(req)
    hashNum := ihash(clientIP)
    bucketNum := hashNum % len(s.bucket)
    currBucketNum := bucketNum

    s.RLock()
    defer s.RUnlock()

    servers := s.bucket[currBucketNum]
    addresses := make([]string, 0, len(servers))
    for len(servers) == 0 {
        currBucketNum++
        if currBucketNum == 10 {
            currBucketNum = 0
        }

        if currBucketNum == bucketNum {
            return "", ErrNoServer
        }
        servers = s.bucket[currBucketNum]
    }

    for address := range servers {
        addresses = append(addresses, address)
    }

    n := len(addresses)
    // TODO: Should have another logic picking the servers from the bucket. Use rand for now.
    index := s.rand.Intn(n)
    randomAddr := addresses[index]

    return randomAddr, nil
}

// Renew updates the bucket within SIH with the given healthyServers.
func (s *SIH) Renew(currentHealthyServers model.BEServers) {
    s.Lock()
    defer s.Unlock()

    // 1. Remove down servers.
    for _, bucket := range s.bucket {
        for serverAddr := range bucket {
            if _, ok := currentHealthyServers[serverAddr]; !ok {
                delete(bucket, serverAddr)
            }
        }
    }
    // 2. Update healthy servers.
    for addr := range currentHealthyServers {
        if _, ok := s.exists(addr); !ok {
            // 1. Hash the address and throw it into the bucket it belongs.
            bucketNum := ihash(addr) % len(s.bucket)
            s.bucket[bucketNum][addr] = struct{}{}
        }
    }
}

// exists check if a serverAddress is in the bucket.
func (s *SIH) exists(serverAddress string) (int, bool) {
    // Don't have to lock here. The methods that calls exists already lock on the outside.
    // Re-locking a locked lock will cause deadlock.
    for bucketNum := range s.bucket {
        if _, ok := s.bucket[bucketNum][serverAddress]; ok {
            return bucketNum, true
        }
    }
    return -1, false
}

// ihash is used as ihash(key) % ip length to choose a bucket number.
func ihash(key string) int {
    h := fnv.New32a()
    _, _ = h.Write([]byte(key))
    return int(h.Sum32() & 0x7fffffff)
}
