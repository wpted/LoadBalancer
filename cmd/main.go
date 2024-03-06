package main

import (
    "LoadBalancer/internal/lb"
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    // scanPeriod is defaulted to 10 seconds.
    scanPeriod := flag.Int("t", 10, "scan period")

    // algoBrief is defaulted to Round-Robin.
    algoBrief := flag.String("algo", "RR", "load balancing algorithm")

    flag.Parse()

    srv, err := lb.New(80, *scanPeriod, *algoBrief)
    if err != nil {
        panic(err)
    }
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
