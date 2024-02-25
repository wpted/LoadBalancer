package main

import (
    "LoadBalancer/internal/lb"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := lb.New(80)
    srv.Start()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    // Wait for os signal to come in.
    <-sigChan
    log.Println("Shutting down load balancer.")

    // Wait for 5 seconds.
    time.Sleep(5 * time.Second)

    // Inform ServerScan to shut down.
    srv.Close()
    log.Println("Load balancer shut down succeeded.")
}
