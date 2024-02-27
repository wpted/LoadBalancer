package lbalgo

import (
    "LoadBalancer/internal/lb"
    "net/http"
    "testing"
)

func TestSRR_ChooseServer(t *testing.T) {
    bes := lb.BEServers{
        "Address A": new(lb.BEServer),
        "Address B": new(lb.BEServer),
        "Address C": new(lb.BEServer),
        "Address D": new(lb.BEServer),
    }

    srr := NewSRR(bes)

    testCases := []struct {
        clientReq      *http.Request
        expectedChosen string
    }{
        {clientReq: &http.Request{RemoteAddr: "10.0.0.1"}},
        {clientReq: &http.Request{RemoteAddr: "10.0.0.2"}},
        {clientReq: &http.Request{RemoteAddr: "10.0.0.3"}},
    }

    for _, tc := range testCases {
        chosen, err := srr.ChooseServer(tc.clientReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }

        // Set the first chosen server to test case.
        tc.expectedChosen = chosen

        // Choose again and compare.
        chosen, err = srr.ChooseServer(tc.clientReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }

        if chosen != tc.expectedChosen {
            t.Errorf("error choosing server: expected %s, got %s.\n", tc.expectedChosen, chosen)
        }
    }
}

func TestSRR_Renew(t *testing.T) {
    bes := lb.BEServers{
        "Address A": new(lb.BEServer),
        "Address B": new(lb.BEServer),
        "Address C": new(lb.BEServer),
        "Address D": new(lb.BEServer),
    }

    srr := NewSRR(bes)

    allClients := Clients{
        "10.0.0.1": "Address A",
        "10.0.0.2": "Address B",
        "10.0.0.3": "Address C",
    }

    srr.AllClients = allClients

    newBes := lb.BEServers{
        "Address B": new(lb.BEServer),
        "Address C": new(lb.BEServer), // Delete server A, D.
        "Address E": new(lb.BEServer), // Add server E.
    }

    srr.Renew(newBes)

    for client, originalBackendServer := range allClients {
        if _, ok := newBes[originalBackendServer]; ok {
            // Healthy server should remain.
            updatedServer := srr.AllClients[client]
            if updatedServer != originalBackendServer {
                t.Errorf("error renewing server: expected %s, got %s.\n", originalBackendServer, updatedServer)
            }
        } else {
            // Unhealthy server should be updated with the next element in queue.
            updatedServer := srr.AllClients[client]
            if updatedServer == originalBackendServer {
                t.Errorf("error renewing server: should be except but %s.\n", originalBackendServer)
            }
        }
    }
}
