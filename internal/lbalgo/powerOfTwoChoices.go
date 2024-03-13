package lbalgo

import (
    "LoadBalancer/internal/model"
    "math/rand"
    "net/http"
    "sync"
    "time"
)

const TWO = 2

// PTC is the struct used for Power of Two Choices.
type PTC struct {
    sync.RWMutex
    servers []model.BEServer
}

// NewPTC creates a PTC instance.
func NewPTC(backendServers *model.BEServers) *PTC {
    servers := make([]model.BEServer, 0)
    if backendServers != nil {
        for _, srv := range *backendServers {
            servers = append(servers, *srv)
        }
    }

    ptc := &PTC{
        servers: servers,
    }

    return ptc
}

// ChooseServer chooses a server based comparison result on the randomly selected two server.
func (p *PTC) ChooseServer(_ *http.Request) (string, error) {
    selected, err := p.choose(TWO)
    if err != nil {
        return "", err
    }

    return chooseLeastConnection(selected), nil
}

// Renew updates the list within PTC with the given healthyServers.
func (p *PTC) Renew(currentHealthyServers model.BEServers) {
    p.Lock()
    defer p.Unlock()

    newServer := make([]model.BEServer, 0)
    // 1. Check down servers.
    for _, srv := range p.servers {
        if healthySrv, ok := currentHealthyServers[srv.Address]; ok {
            // Copy the healthy server.
            newServer = append(newServer, *healthySrv)
        }
    }

    // 2. Check up servers.
    for addr, newSrv := range currentHealthyServers {
        if !p.exists(addr) {
            // Copy the new server.
            newServer = append(newServer, *newSrv)
        }
    }

    p.servers = newServer
}

// exists check if a serverAddress is in the list.
func (p *PTC) exists(address string) bool {
    for _, srv := range p.servers {
        if srv.Address == address {
            return true
        }
    }

    return false
}

// chooseLeastConnection selects a server with the least connections.
func chooseLeastConnection(servers []model.BEServer) string {
    leastConnectionServer := servers[0]
    for i := 1; i < len(servers); i++ {
        if servers[i].Connections < leastConnectionServer.Connections {
            leastConnectionServer = servers[i]
        }
    }

    return leastConnectionServer.Address
}

// choose selects k servers from PTC randomly.
// If the length of p.Servers are smaller than k, return all objects that exists.
// Returns an error if there's no server in p.
func (p *PTC) choose(k int) ([]model.BEServer, error) {
    p.Lock()
    defer p.Unlock()
    if len(p.servers) == 0 {
        return nil, ErrNoServer
    }

    if k > len(p.servers) {
        return p.servers, nil
    }

    // Seed the source with now.
    rand.New(rand.NewSource(time.Now().UnixNano()))

    // Shuffle the list.
    rand.Shuffle(len(p.servers), func(i, j int) {
        p.servers[i], p.servers[j] = p.servers[j], p.servers[i]
    })

    return p.servers[:k], nil
}
