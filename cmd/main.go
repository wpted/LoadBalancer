package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    lb := New()
    lb.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        log.Println(req)
        backendResponse := lb.Forward(req)
        // Write response back to client.
        _, _ = fmt.Fprint(w, backendResponse.Body)
    })

    err := http.ListenAndServe(":80", lb)
    if err != nil {
        panic(err)
    }
}

type LoadBalancer struct {
    http.ServeMux
    http.Client
    BackendServers []BackendServer
}

func New() *LoadBalancer {
    return &LoadBalancer{
        BackendServers: []BackendServer{
            {Address: "http://localhost:1080", Alive: true},
        },
    }
}

type BackendServer struct {
    Address string
    Alive   bool
}

func (l *LoadBalancer) Forward(request *http.Request) *http.Response {
    // 1. Forward the request to a address from the Server lists.
    addr := l.BackendServers[0]

    r, err := http.NewRequest(request.Method, addr.Address, request.Body)
    if err != nil {
        log.Println(err)
    }

    // Response from backend service.
    resp, err := l.Do(r)
    if err != nil {
        log.Println(err)
    }
    log.Println(resp)

    // Return the response

    return resp
}
