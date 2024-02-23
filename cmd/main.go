package main

import (
    "LoadBalancer/internal/lb"
    "net/http"
)

func main() {
    srv := lb.New()
    srv.HandleFunc("/", srv.Forward)
    srv.HandleFunc("/register", srv.Register)

    err := http.ListenAndServe(":80", srv)
    if err != nil {
        panic(err)
    }
}
