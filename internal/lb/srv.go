package lb

import (
    "net/http"
    "sync"
)

// LoadBalancer distributes traffic to AliveServers.
type LoadBalancer struct {
    http.ServeMux
    http.Client
    sync.RWMutex
    AliveServers map[string]struct{}
    DownServers  map[string]struct{}
}

// New creates an instance of LoadBalancer.
func New() *LoadBalancer {
    return &LoadBalancer{
        AliveServers: make(map[string]struct{}),
        DownServers:  make(map[string]struct{}),
    }
}
