package lb

import (
    "encoding/json"
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

// RegisterRequest is used for registering backend servers.
type RegisterRequest struct {
    Address string `json:"address"`
}

// Register is a handler that is used by endpoint '/register'.
func (l *LoadBalancer) Register(w http.ResponseWriter, req *http.Request) {
    var p RegisterRequest
    decoder := json.NewDecoder(req.Body)
    if err := decoder.Decode(&p); err != nil {
        // Return error message.
    }
    // Ping the address.
    serverAlive := l.healthCheck(p.Address)

    l.RLock()
    defer l.RUnlock()
    // Only register server when backend server is alive.
    if serverAlive {
        l.AliveServers[p.Address] = struct{}{}
    }
}
