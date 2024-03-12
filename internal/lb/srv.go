package lb

import (
    "LoadBalancer/internal/lb/response"
    "LoadBalancer/internal/lbalgo"
    "LoadBalancer/internal/model"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "sync"
    "time"
)

// LoadBalancer distributes traffic to AliveServers.
type LoadBalancer struct {
    http.ServeMux
    http.Client
    sync.RWMutex
    Port         int
    AliveServers model.BEServers
    DownServers  model.BEServers
    ScanDone     chan struct{}
    ScanPeriod   time.Duration
    AlgoDriver   lbalgo.LBAlgo
}

// New creates an instance of LoadBalancer.
func New(port int, scanPeriod int, algoBrief string) (*LoadBalancer, error) {
    algoDriver, err := lbalgo.ChooseAlgo(algoBrief)
    if err != nil {
        return nil, err
    }

    return &LoadBalancer{
        Port:         port,
        AliveServers: make(map[string]*model.BEServer),
        DownServers:  make(map[string]*model.BEServer),
        ScanDone:     make(chan struct{}),
        ScanPeriod:   time.Duration(scanPeriod) * time.Second,
        AlgoDriver:   algoDriver, // no server in the algo driver now.
    }, nil
}

// Start starts the server.
// The method spawns two goroutines, one starting the http server and the other start the periodic scan routine.
func (l *LoadBalancer) Start() {
    l.HandleFunc("/", l.Forward)
    l.HandleFunc("/register", l.Register)

    go func() {
        if err := http.ListenAndServe(fmt.Sprintf(":%d", l.Port), l); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("Load balancer server error: %v", err)
        }
        return
    }()

    go l.ScanPeriodically()
}

// Close shuts down all goroutines and closes the Done channel.
func (l *LoadBalancer) Close() {
    // Send two signal to the done channel.
    // This shuts down ScanPeriodically().
    l.ScanDone <- struct{}{}
    close(l.ScanDone)
}

// RegisterRequest is used for registering backend servers.
type RegisterRequest struct {
    Address string `json:"address"`
    Weight  int    `json:"weight"`
}

// Register is a handler that is used by endpoint '/register'.
func (l *LoadBalancer) Register(w http.ResponseWriter, req *http.Request) {
    var p RegisterRequest
    decoder := json.NewDecoder(req.Body)
    decoder.DisallowUnknownFields()
    if err := decoder.Decode(&p); err != nil {
        // Return error message.
        response.WriteJsonResponse(w, http.StatusInternalServerError, response.NewErrorResponse(err))
        return
    }
    // Ping the address.
    serverAlive := l.healthCheck(p.Address)
    l.RLock()
    defer l.RUnlock()
    // Only register server when backend server is alive.
    if serverAlive {
        l.AliveServers[p.Address] = model.NewBEServer(p.Address, p.Weight)
        responsePayload := response.NewSuccessResponse(
            struct {
                Server string `json:"server"`
                Weight int    `json:"weight"`
            }{
                Server: p.Address,
                Weight: p.Weight,
            })
        response.WriteJsonResponse(w, http.StatusOK, responsePayload)
        return
    } else {
        // Return service not alive, registration failed.
        responsePayload := response.NewFailResponse(
            struct {
                Title string `json:"title"`
            }{Title: fmt.Sprintf("%s not alive, registration failed.", p.Address)})
        response.WriteJsonResponse(w, http.StatusNotFound, responsePayload)
        return
    }
}

// Forward is a handler that distributes traffic to all AliveServers.
func (l *LoadBalancer) Forward(w http.ResponseWriter, req *http.Request) {
    // 1. Forward the request to an address from the Server lists.
    addr, err := l.AlgoDriver.ChooseServer(req)
    if err != nil {
        log.Println(err)
        response.WriteJsonResponse(w, http.StatusServiceUnavailable, response.NewErrorResponse(err))
        return
    }

    newReq, err := copyRequest(req, addr)

    if err != nil {
        log.Println(err)
    }

    // Response from backend service.
    resp, err := l.Do(newReq)
    if err != nil {
        log.Println(err)
    }
    defer func() {
        err = resp.Body.Close()
        if err != nil {
            log.Fatal(err)
        }
    }()

    bodyBytes, err := io.ReadAll(resp.Body)
    // Write response back to client.
    _, err = fmt.Fprint(w, fmt.Sprintf("From backend server: %s, data: [ '%s' ].\n", addr, string(bodyBytes)))
    if err != nil {
        log.Println(err)
    }

}

func copyRequest(req *http.Request, target string) (*http.Request, error) {
    // The general form represented is: [scheme:][//[userinfo@]host][/]path[?query][#fragment]
    r, err := http.NewRequest(req.Method, target, req.Body)
    if err != nil {
        return nil, err
    }

    // Deep copy the header instead of using the original one
    r.Header = make(http.Header)
    for k, v := range req.Header {
        r.Header[k] = v
    }
    return r, nil
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

// ScanPeriodically triggers the scan periodically in a different goroutine.
func (l *LoadBalancer) ScanPeriodically() {
    scanTicker := time.NewTicker(l.ScanPeriod)
    defer scanTicker.Stop()
    for {
        select {
        case <-l.ScanDone:
            return
        case <-scanTicker.C:
            l.scanServers()
        }
    }
}

// scanServers checks all registered servers.
// This method enables the load balancer to manage servers that come back online after passing health checks and to remove servers that failed.
func (l *LoadBalancer) scanServers() {
    l.RLock()
    // Check all servers in AliveServers.
    for addr := range l.AliveServers {
        healthy := l.healthCheck(addr)
        if !healthy {
            l.DownServers[addr] = l.AliveServers[addr]
            delete(l.AliveServers, addr)
        }
    }

    // Check all servers in DownServers.
    for addr := range l.DownServers {
        healthy := l.healthCheck(addr)
        if healthy {
            l.AliveServers[addr] = l.DownServers[addr]
            delete(l.DownServers, addr)
        }
    }
    // Update the algo driver with current servers that are alive.
    l.AlgoDriver.Renew(l.AliveServers)

    l.RUnlock()
}
