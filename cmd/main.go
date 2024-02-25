package main

import (
    "LoadBalancer/internal/lb"
    "errors"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := lb.New()
    srv.HandleFunc("/", srv.Forward)
    srv.HandleFunc("/register", srv.Register)

    go func() {
        if err := http.ListenAndServe(":80", srv); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("Load balancer server error: %v", err)
        }
    }()

    go srv.ServerScan()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    // Wait for os signal to come in.
    <-sigChan
    log.Println("Shutting down load balancer.")

    // Wait for 5 seconds.
    time.Sleep(5 * time.Second)

    srv.Close()
    log.Println("Load balancer shut down succeeded.")
}
