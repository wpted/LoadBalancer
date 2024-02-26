package lb

import (
    "encoding/json"
    "errors"
    "fmt"
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
    ScanOptions
    Port         int
    AliveServers map[string]struct{}
    DownServers  map[string]struct{}
    ReplDone     chan struct{}
}

// New creates an instance of LoadBalancer.
func New(port int) *LoadBalancer {
    scanOpt := ScanOptions{
        ScanDone:   make(chan struct{}),
        ScanPeriod: 10 * time.Second, // Scan Period default to 10 seconds.
        ScanTicker: time.NewTicker(10 * time.Second),
    }
    return &LoadBalancer{
        ScanOptions:  scanOpt,
        Port:         port,
        AliveServers: make(map[string]struct{}),
        DownServers:  make(map[string]struct{}),
        //ReplDone:     make(chan struct{}),
    }
}

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
    //go l.Repl()
}

// Close shuts down all goroutines and closes the Done channel.
func (l *LoadBalancer) Close() {
    // Send two signal to the done channel.
    // This shuts down Repl() and ScanPeriodically().
    l.ScanDone <- struct{}{}
    //l.ReplDone <- struct{}{}
    close(l.ScanDone)
    //close(l.ReplDone)
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

    // 1. Forward the request to an address from the Server lists.
    addr := ""
    newReq, err := copyRequest(req, addr)

    if err != nil {
        log.Println(err)
    }

    // Response from backend service.
    resp, err := l.Do(newReq)
    if err != nil {
        log.Println(err)
    }
    log.Println(resp)

    // Write response back to client.
    _, _ = fmt.Fprint(w, resp.Body)
}

func copyRequest(req *http.Request, target string) (*http.Request, error) {
    // The general form represented is: [scheme:][//[userinfo@]host][/]path[?query][#fragment]
    fmt.Printf("%#v", req.URL.Scheme)
    r, err := http.NewRequest(req.Method, target, req.Body)
    if err != nil {
        return nil, err
    }

    // Copy the detailed path.
    r.URL.Scheme = req.URL.Scheme
    r.URL.Path = req.URL.Path

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

type ScanOptions struct {
    ScanDone   chan struct{}
    ScanPeriod time.Duration
    ScanTicker *time.Ticker
}

func (l *LoadBalancer) ScanPeriodically() {
    for {
        select {
        case <-l.ScanDone:
            l.ScanTicker.Stop()
            return
        case <-l.ScanTicker.C:
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
            delete(l.AliveServers, addr)
            l.DownServers[addr] = struct{}{}
        }
    }

    // Check all servers in DownServers.
    for addr := range l.DownServers {
        healthy := l.healthCheck(addr)
        if healthy {
            delete(l.DownServers, addr)
            l.AliveServers[addr] = struct{}{}
        }
    }
    l.RUnlock()
}

//func (l *LoadBalancer) Repl() {
//    // By default, we start a repl waiting for changes on the scan period.
//    for {
//        select {
//        case <-l.ReplDone:
//            return
//        default:
//            fmt.Print("Load balancer repl: ")
//            var input string
//            _, err := fmt.Scanf("%s", &input)
//            if err != nil {
//                log.Println(err)
//            }
//            if strings.ToLower(input) == "scan" {
//                var scanPeriod int
//                fmt.Print("New Scan period: ")
//                _, err = fmt.Scanf("%d", &scanPeriod)
//                if err != nil {
//                    log.Println(err)
//                }
//
//                t := time.Duration(scanPeriod) * time.Second
//
//                if t != l.ScanPeriod {
//                    // Reset the ticker and update l.ScanPeriod.
//                    l.ScanPeriod = t
//                    l.ScanTicker.Reset(l.ScanPeriod)
//                }
//            }
//        }
//    }
//}
