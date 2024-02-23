package lb

import (
    "encoding/json"
    "fmt"
    "log"
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

// Forward is a handler that distributes traffic to all AliveServers.
func (l *LoadBalancer) Forward(w http.ResponseWriter, req *http.Request) {
    log.Println(req)

    // 1. Forward the request to an address from the Server lists.
    addr := ""

    r, err := http.NewRequest(req.Method, addr, req.Body)
    if err != nil {
        log.Println(err)
    }

    // Response from backend service.
    resp, err := l.Do(r)
    if err != nil {
        log.Println(err)
    }
    log.Println(resp)

    // Write response back to client.
    _, _ = fmt.Fprint(w, resp.Body)
}

// healthCheck sends a request to the targetServer.
// Returns a boolean representing the server health status.
func (l *LoadBalancer) healthCheck(targetServer string) bool {
    healthCheckEndpoint := fmt.Sprintf("%s/health", targetServer)
    resp, err := http.Get(healthCheckEndpoint)
    if err != nil {
        return false
    }

    if resp.StatusCode != http.StatusOK {
        return false
    }

    return true
}
