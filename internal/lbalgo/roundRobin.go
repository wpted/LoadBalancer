package lbalgo

import (
    "LoadBalancer/internal/lb"
    "net/http"
    "sort"
    "sync"
)

// RR is the implemented queue for using the Round Robin algorithm.
type RR struct {
    sync.RWMutex
    servers []string
}

func (r *RR) Len() int      { return len(r.servers) }
func (r *RR) Swap(i, j int) { r.servers[i], r.servers[j] = r.servers[j], r.servers[i] }
func (r *RR) Less(i, j int) bool {
    return r.servers[i] < r.servers[j]
}

// NewRR creates a new instance of RR.
func NewRR(backendServers lb.BEServers) *RR {
    servers := make([]string, 0)
    for addr := range backendServers {
        servers = append(servers, addr)
    }

    rr := &RR{servers: servers}
    // Since the order isn't consistent when reading from a map, sort the result.
    sort.Sort(rr)
    return rr
}

// Renew updates the queue within RR.
func (r *RR) Renew(backendServers map[string]struct{}) {
    // 1. Check down servers.
    for _, addr := range r.servers {
        if _, ok := backendServers[addr]; !ok {
            // Means that there's server down.
            r.remove(addr)
        }
    }

    // 2. Check up servers.
    for addr := range backendServers {
        if !r.exists(addr) {
            r.push(addr)
        }
    }
}

// ChooseServer rotates the queue within RR and returns the chosenServer.
func (r *RR) ChooseServer(_ *http.Request) (string, error) {
    chosenServer := r.rotate()
    if chosenServer == "" {
        return "", ErrNoServer
    }

    return chosenServer, nil
}

// rotate rotates the queue within RR.
func (r *RR) rotate() string {
    head := r.pop()
    if head != "" {
        r.push(head)
    }
    return head
}

// push adds address to the end of the queue within RR.
func (r *RR) push(serverAddress string) {
    r.Lock()
    defer r.Unlock()

    r.servers = append(r.servers, serverAddress)
}

// pop removes the first element from the queue within RR.
func (r *RR) pop() string {
    r.Lock()
    defer r.Unlock()

    var head string
    if len(r.servers) != 0 {
        head = r.servers[0]
        r.servers = r.servers[1:]
    }
    return head
}

// exists checks whether a serverAddress exists within RR.
func (r *RR) exists(serverAddress string) bool {
    r.RLock()
    defer r.RUnlock()
    // TODO: Should use binary search.
    for _, addr := range r.servers {
        if serverAddress == addr {
            return true
        }
    }

    return false
}

// remove removes the serverAddress from the queue within RR.
func (r *RR) remove(serverAddress string) {
    r.Lock()
    defer r.Unlock()
    for n, addr := range r.servers {
        if serverAddress == addr {
            r.servers = append(r.servers[:n], r.servers[n+1:]...)
            return
        }
    }
}
