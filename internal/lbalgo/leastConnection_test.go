package lbalgo

import (
    "LoadBalancer/internal/lb"
    "net/http"
    "testing"
)

func TestLC_ChooseServer(t *testing.T) {
    testCases := []struct {
        servers        *lb.BEServers
        expectedChosen string
    }{
        {
            servers: &lb.BEServers{
                "Address A": &lb.BEServer{Address: "Address A", Connections: 4},
                "Address B": &lb.BEServer{Address: "Address B", Connections: 1},
                "Address C": &lb.BEServer{Address: "Address C", Connections: 3},
                "Address D": &lb.BEServer{Address: "Address D", Connections: 2},
            },
            expectedChosen: "Address B",
        },
        {
            servers: &lb.BEServers{
                "Address A": &lb.BEServer{Address: "Address A", Connections: 4},
                "Address B": &lb.BEServer{Address: "Address B", Connections: 1},
                "Address C": &lb.BEServer{Address: "Address C", Connections: 3},
                "Address D": &lb.BEServer{Address: "Address D", Connections: 2},
                "Address E": &lb.BEServer{Address: "Address E", Connections: 5},
                "Address G": &lb.BEServer{Address: "Address G", Connections: 8},
                "Address F": &lb.BEServer{Address: "Address F", Connections: 12},
                "Address p": &lb.BEServer{Address: "Address P", Connections: 6},
            },
            expectedChosen: "Address B",
        },
    }

    emptyReq := new(http.Request)
    for _, tc := range testCases {
        lc := NewLC(tc.servers)
        chosen, err := lc.ChooseServer(emptyReq)
        if err != nil {
            t.Errorf("error choosing server: %v.\n", err)
        }

        if chosen != tc.expectedChosen {
            t.Errorf("error choosing server: expected %s, got %s.\n", tc.expectedChosen, chosen)
        }
    }
}

func TestLC_Renew(t *testing.T) {
    bes := lb.BEServers{
        "Address A": &lb.BEServer{Address: "Address A", Connections: 1},
        "Address B": &lb.BEServer{Address: "Address B", Connections: 2},
        "Address C": &lb.BEServer{Address: "Address C", Connections: 3},
        "Address D": &lb.BEServer{Address: "Address D", Connections: 4},
    }

    testCases := []struct {
        newBes         lb.BEServers
        expectedChosen string
    }{
        {
            newBes: lb.BEServers{
                "Address B": &lb.BEServer{Address: "Address B", Connections: 2},
                "Address C": &lb.BEServer{Address: "Address C", Connections: 3},
                "Address D": &lb.BEServer{Address: "Address D", Connections: 4},
                "Address A": &lb.BEServer{Address: "Address A", Connections: 1},
                "Address E": &lb.BEServer{Address: "Address E", Connections: 6}, // Add server E.
            },
            expectedChosen: "Address A",
        },
        {
            newBes: lb.BEServers{
                "Address B": &lb.BEServer{Address: "Address B", Connections: 2},
                "Address C": &lb.BEServer{Address: "Address C", Connections: 3},
                "Address D": &lb.BEServer{Address: "Address D", Connections: 4}, // Delete server A.
            },
            expectedChosen: "Address B",
        },
        {
            newBes: lb.BEServers{
                "Address B": &lb.BEServer{Address: "Address B", Connections: 2},
                "Address E": &lb.BEServer{Address: "Address E", Connections: 4}, // Delete server A, C, D and add server E.
            },
            expectedChosen: "Address B",
        },
    }

    emptyReq := new(http.Request)
    for _, tc := range testCases {
        lc := NewLC(&bes)
        lc.Renew(tc.newBes)
        chosen, err := lc.ChooseServer(emptyReq)
        if err != nil {
            t.Errorf("error choosing server: %v.\n", err)
        }

        if chosen != tc.expectedChosen {
            t.Errorf("error choosing server: expected %s, got %s.\n", tc.expectedChosen, chosen)
        }
    }
}

func Test_buildMinHeap(t *testing.T) {
    testCases := []struct {
        servers         *lb.BEServers
        expectedMinHeap []*lb.BEServer
    }{
        {
            servers: &lb.BEServers{
                "Address A": &lb.BEServer{Connections: 4},
                "Address B": &lb.BEServer{Connections: 1},
                "Address C": &lb.BEServer{Connections: 3},
                "Address D": &lb.BEServer{Connections: 2},
            },
        },
        {
            servers: &lb.BEServers{
                "Address A": &lb.BEServer{Connections: 4},
                "Address B": &lb.BEServer{Connections: 1},
                "Address C": &lb.BEServer{Connections: 3},
                "Address E": &lb.BEServer{Connections: 2},
                "Address F": &lb.BEServer{Connections: 6},
                "Address G": &lb.BEServer{Connections: 10},
                "Address H": &lb.BEServer{Connections: 11},
            },
        },
    }

    for _, tc := range testCases {
        lc := NewLC(tc.servers)
        lc.buildMinHeap()

        if !isMinHeap(lc.servers) {
            t.Errorf("Error build minHeap")
        }
    }
}

func isMinHeap(servers []*lb.BEServer) bool {
    // A min heap should satisfy condition: Every node in array must be smaller than its children ( If exist ).
    l := len(servers)
    for n, srv := range servers {
        lChildIdx := leftChildIdx(n)
        rChildIdx := rightChildIdx(n)

        if lChildIdx < l && srv.Connections > servers[lChildIdx].Connections {
            return false
        }

        if rChildIdx < l && srv.Connections > servers[rChildIdx].Connections {
            return false
        }
    }
    return true
}
