package lbalgo

import (
    "LoadBalancer/internal/model"
    "net/http"
    "testing"
)

func TestWRR_ChooseServer(t *testing.T) {
    bes := model.BEServers{
        "Address A": &model.BEServer{Weight: 2},
        "Address B": &model.BEServer{Weight: 1},
        "Address C": &model.BEServer{Weight: 5},
    }

    testCases := []struct {
        requests       int
        expectedChosen string
    }{
        {
            requests:       3,
            expectedChosen: "Address C",
        },
        {
            requests:       5,
            expectedChosen: "Address C",
        },
        {
            requests:       6,
            expectedChosen: "Address A",
        },
        {
            requests:       8,
            expectedChosen: "Address B",
        },
        {
            requests:       9,
            expectedChosen: "Address C",
        },
        {
            requests:       14,
            expectedChosen: "Address A",
        },
    }

    for _, tc := range testCases {
        wrr := NewWRR(&bes)

        var chosen string
        var err error
        emptyReq := new(http.Request)
        for i := 0; i < tc.requests; i++ {
            chosen, err = wrr.ChooseServer(emptyReq)
            if err != nil {
                t.Errorf("error choosing server: got %#v.\n", err)
            }
        }

        if chosen != tc.expectedChosen {
            t.Errorf("error choosing server: expected %s, got %s.\n", tc.expectedChosen, chosen)
        }
    }
}

func TestWRR_Renew(t *testing.T) {
    bes := model.BEServers{
        "Address A": &model.BEServer{Weight: 1},
        "Address B": &model.BEServer{Weight: 2},
        "Address C": &model.BEServer{Weight: 3},
        "Address D": &model.BEServer{Weight: 4},
    }

    testCases := []struct {
        newBes   model.BEServers
        expected []string
    }{
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{Weight: 2},
                "Address C": &model.BEServer{Weight: 3},
                "Address D": &model.BEServer{Weight: 4},
                "Address A": &model.BEServer{Weight: 1},
                "Address E": &model.BEServer{Weight: 6}, // Add server E.
            },
            expected: []string{"Address E", "Address D", "Address C", "Address B", "Address A"},
        },
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{Weight: 2},
                "Address C": &model.BEServer{Weight: 3},
                "Address D": &model.BEServer{Weight: 4}, // Delete server A.
            },
            expected: []string{"Address D", "Address C", "Address B"},
        },
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{Weight: 2},
                "Address E": &model.BEServer{Weight: 4}, // Delete server A, C, D and add server E.
            },
            expected: []string{"Address E", "Address B"},
        },
    }

    for _, tc := range testCases {
        wrr := NewWRR(&bes)
        wrr.Renew(tc.newBes)
        serverAddresses := make([]string, 0)
        for _, server := range wrr.servers {
            serverAddresses = append(serverAddresses, server.Addr)
        }

        if !assertEqualSlice(serverAddresses, tc.expected) {
            t.Errorf("error renewing server: expected %#v, got %#v.\n", tc.expected, serverAddresses)
        }
    }
}
