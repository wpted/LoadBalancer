package model

import "time"

type BEServers map[string]*BEServer
type BEServer struct {
    Address        string
    Weight         int
    ConnectionTime time.Duration
    Connections    int
}

// NewBEServer creates a new instance of BEServer.
func NewBEServer(address string, weight int) *BEServer {
    return &BEServer{
        Address: address,
        Weight:  weight,
    }
}
