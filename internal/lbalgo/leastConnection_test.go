package lbalgo

import (
    "LoadBalancer/internal/lb"
    "testing"
)

func TestLC_ChooseServer(t *testing.T) {

}

func TestLC_Renew(t *testing.T) {

}

func Test_buildMinHeap(t *testing.T) {

    testCases := []struct {
        servers         *lb.BEServers
        expectedMinHeap []*lb.BEServer
    }{
        {
            servers: &lb.BEServers{
                "Address A": &lb.BEServer{Weight: 4},
                "Address B": &lb.BEServer{Weight: 1},
                "Address C": &lb.BEServer{Weight: 3},
                "Address D": &lb.BEServer{Weight: 2},
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
