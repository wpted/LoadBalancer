package lbalgo

import (
    "LoadBalancer/internal/model"
    "net/http"
    "sort"
    "sync"
)

type WRR struct {
    sync.RWMutex
    servers []*weightedServer
}

type weightedServer struct { // We need an extra data structure to keep track of the Count.
    Addr   string
    Weight int
    Count  int
}

func (w *WRR) Len() int      { return len(w.servers) }
func (w *WRR) Swap(i, j int) { w.servers[i], w.servers[j] = w.servers[j], w.servers[i] }
func (w *WRR) Less(i, j int) bool {
    return w.servers[i].Weight > w.servers[j].Weight
}

func NewWRR(backendServers *model.BEServers) *WRR {
    servers := make([]*weightedServer, 0)
    if backendServers != nil {
        for addr, srv := range *backendServers {
            ws := &weightedServer{
                Addr:   addr,
                Weight: srv.Weight,
                Count:  srv.Weight,
            }
            servers = append(servers, ws)
        }
    }
    wrr := &WRR{servers: servers}
    // Sort by weight.
    sort.Sort(wrr)
    return wrr
}

func (w *WRR) ChooseServer(_ *http.Request) (string, error) {
    if len(w.servers) != 0 {
        chosenServer := w.servers[0].Addr
        w.servers[0].Count--
        if w.servers[0].Count <= 0 {
            w.servers[0].Count = w.servers[0].Weight
            w.rotate()
        }
        return chosenServer, nil
    }
    return "", ErrNoServer
}

func (w *WRR) Renew(currentHealthyServers model.BEServers) {
    // 1. Check down servers.
    for _, server := range w.servers {
        if _, ok := currentHealthyServers[server.Addr]; !ok {
            // Means that there's server down.
            w.remove(server.Addr)
        }
    }

    // 2. Check up servers.
    for addr, server := range currentHealthyServers {
        if srv, exist := w.exists(addr); !exist {
            ws := &weightedServer{
                Addr:   addr,
                Weight: server.Weight,
                Count:  server.Weight,
            }
            w.push(ws)
        } else {
            srv.Weight = server.Weight
            srv.Count = server.Weight
        }
    }

    sort.Sort(w)
}

func (w *WRR) rotate() *weightedServer {
    head := w.pop()
    if head != nil {
        w.push(head)
    }
    return head
}

func (w *WRR) push(server *weightedServer) {
    w.Lock()
    defer w.Unlock()

    w.servers = append(w.servers, server)
}

func (w *WRR) pop() *weightedServer {
    w.Lock()
    defer w.Unlock()

    var head *weightedServer
    if len(w.servers) != 0 {
        head = w.servers[0]
        w.servers = w.servers[1:]
    }
    return head
}

func (w *WRR) exists(serverAddress string) (*weightedServer, bool) {
    w.RLock()
    defer w.RUnlock()
    // TODO: Should use binary search.
    for _, server := range w.servers {
        if serverAddress == server.Addr {
            return server, true
        }
    }
    return nil, false
}

func (w *WRR) remove(serverAddress string) {
    w.Lock()
    defer w.Unlock()

    newServers := make([]*weightedServer, 0)
    for _, server := range w.servers {
        if serverAddress != server.Addr {
            newServers = append(newServers, server)
        }
    }
    w.servers = newServers
}
