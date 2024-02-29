package lbalgo

import (
    "LoadBalancer/internal/lb"
    "net/http"
    "sync"
)

// Using Go's sort interface would be easier, but for practice I'm implementing a minimum priority queue.

type LC struct {
    sync.RWMutex
    servers []*lb.BEServer // Using a slice as a binary heap.
}

func NewLC(backendServers *lb.BEServers) *LC {
    servers := make([]*lb.BEServer, 0)
    for _, srv := range *backendServers {
        servers = append(servers, srv)
    }
    return &LC{servers: servers}
}

func (l *LC) ChooseServer(_ *http.Request) (string, error) {
    // Call buildMinHeap().
    if len(l.servers) != 0 {
        l.buildMinHeap()
        return l.servers[0].Address, nil
    }

    return "", ErrNoServer
}

func (l *LC) Renew(backendServers lb.BEServers) {
    // 1. Check down servers.
    for _, srv := range l.servers {
        if _, ok := backendServers[srv.Address]; !ok {
            // Means that there's server down.
            l.remove(srv.Address)
        }
    }

    // 2. Check up servers.
    for addr, srv := range backendServers {
        if !l.exists(addr) {
            l.push(srv)
        }
    }

    l.buildMinHeap()
}

// exists checks whether a serverAddress exists within LC.
func (l *LC) exists(serverAddress string) bool {
    l.RLock()
    defer l.RUnlock()
    // TODO: Should use binary search.
    for _, srv := range l.servers {
        if serverAddress == srv.Address {
            return true
        }
    }

    return false
}

func (l *LC) push(server *lb.BEServer) {
    l.Lock()
    defer l.Unlock()
    l.servers = append(l.servers, server)
}

func (l *LC) remove(serverAddress string) {
    l.Lock()
    defer l.Unlock()
    newServers := make([]*lb.BEServer, 0)

    for _, srv := range l.servers {
        if serverAddress != srv.Address {
            newServers = append(newServers, srv)
        }
    }

    l.servers = newServers
}

func leftChildIdx(idx int) int {
    return idx*2 + 1
}

func rightChildIdx(idx int) int {
    return idx*2 + 2
}

// minHeapify starts the minimum heapify process from the given index.
func (l *LC) minHeapify(idx int) {
    lowest := idx

    lChildIdx := leftChildIdx(idx)
    rChildIdx := rightChildIdx(idx)

    if lChildIdx < len(l.servers) && l.servers[lChildIdx].Connections < l.servers[lowest].Connections {
        lowest = lChildIdx
    }
    if rChildIdx < len(l.servers) && l.servers[rChildIdx].Connections < l.servers[lowest].Connections {
        lowest = rChildIdx
    }

    if lowest != idx {
        l.servers[lowest], l.servers[idx] = l.servers[idx], l.servers[lowest]
        l.minHeapify(lowest)
    }
}

// buildMinHeap turns l.server into a minimum heap.
func (l *LC) buildMinHeap() {
    l.Lock()
    defer l.Unlock()

    for i := len(l.servers) / 2; i >= 0; i-- {
        l.minHeapify(i)
    }
}
