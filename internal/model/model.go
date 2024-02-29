package model

import "time"

type BEServers map[string]*BEServer
type BEServer struct {
    Address        string
    Weight         int
    ConnectionTime time.Duration
    Connections    int
}
