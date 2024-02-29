package lbalgo

import (
    "LoadBalancer/internal/model"
    "net/http"
    "testing"
)

func TestSRR_ChooseServer(t *testing.T) {
    bes := &model.BEServers{
        "Address A": new(model.BEServer),
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer),
        "Address D": new(model.BEServer),
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
    bes := &model.BEServers{
        "Address A": new(model.BEServer),
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer),
        "Address D": new(model.BEServer),
    }

    srr := NewSRR(bes)

    allClients := Clients{
        "10.0.0.1": "Address A",
        "10.0.0.2": "Address B",
        "10.0.0.3": "Address C",
    }

    srr.AllClients = allClients

    newBes := model.BEServers{
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer), // Delete server A, D.
        "Address E": new(model.BEServer), // Add server E.
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
