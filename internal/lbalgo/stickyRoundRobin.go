package lbalgo

import (
    "LoadBalancer/internal/model"
    "fmt"
    "net/http"
    "sync"
)

type Clients map[string]string // client-ip: server-ip

// SRR instance.
type SRR struct {
    AllClients Clients
    sync.Mutex
    rr *RR
}

func NewSRR(backendServers *model.BEServers) *SRR {
    return &SRR{
        AllClients: make(Clients),
        rr:         NewRR(backendServers),
    }
}

// ChooseServer chooses a backend server for a incoming client.
// It ensures that each client is consistently routed to the same backend server as long as its sticky criteria (IP address) remains the same, providing session affinity or sticky sessions.
func (s *SRR) ChooseServer(req *http.Request) (string, error) {
    clientIP := getClientIP(req)
    fmt.Println(clientIP)
    s.Lock()
    defer s.Unlock()
    beAddr, ok := s.AllClients[clientIP]
    if !ok {
        assignedAddr, err := s.rr.ChooseServer(req)
        if err != nil {
            // Error occurs when there's no server in pool.
            return "", err
        }

        // Store assigned addr.
        s.AllClients[clientIP] = assignedAddr
        return assignedAddr, nil
    }

    return beAddr, nil
}

// Renew updates the round-robin queue and the server bound to the clients.
func (s *SRR) Renew(healthyServers model.BEServers) {
    // Update round-robin queue.
    s.rr.Renew(healthyServers)

    // Go through all clients and replace the unhealthy servers.
    s.Lock()
    defer s.Unlock()

    for client, BEServer := range s.AllClients {
        if _, ok := healthyServers[BEServer]; !ok {
            // Server bound with client not healthy.
            // Replace with new.
            emptyReq := new(http.Request)
            emptyReq.RemoteAddr = client
            newChosenSrv, err := s.rr.ChooseServer(emptyReq)
            if err != nil {
                break // no servers left to assign => empty queue.
            }
            s.AllClients[client] = newChosenSrv
        }
    }
}

// getClientIP gets the IP of the client. If the client is hided behind proxies or load balancers,
// we get the IP from retrieving the value from X-Forwarded-For header.
// This method is only a demo and shouldn't be used in any production code. Header can be changed after a request is sent.
// The best way is to store all occurring ip for further analytics.
func getClientIP(req *http.Request) string {
    clientIP := req.Header.Get("X-Forwarded-For")
    if clientIP == "" {
        return req.RemoteAddr
    }

    return clientIP
}
