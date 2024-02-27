package lbalgo

import (
    "LoadBalancer/internal/lb"
    "errors"
    "net/http"
    "testing"
)

func TestRR_ChooseServer(t *testing.T) {

    emptyReq := new(http.Request)
    t.Run("No servers", func(t *testing.T) {
        bes := new(lb.BEServers)
        rr := NewRR(*bes)
        _, err := rr.ChooseServer(emptyReq)
        if !errors.Is(err, ErrNoServer) {
            t.Errorf("error incorrect error: expected %#v, got %#v.\n", ErrNoServer, err)
        }
    })

    t.Run("Correct Backend Server", func(t *testing.T) {
        bes := lb.BEServers{
            "Address A": new(lb.BEServer),
            "Address B": new(lb.BEServer),
            "Address C": new(lb.BEServer),
            "Address D": new(lb.BEServer),
        }

        // After 1 rotate.
        expectedChosen := "Address A"
        expectedRR := []string{"Address B", "Address C", "Address D", "Address A"}

        rr := NewRR(bes)
        res, err := rr.ChooseServer(emptyReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }

        if res != expectedChosen {
            t.Errorf("error choosing server: expected %s, got %s.\n", expectedChosen, res)
        }

        if !assertEqualSlice(expectedRR, rr.servers) {
            t.Errorf("error choosing server: expected %#v, got %#v.\n", expectedRR, rr.servers)
        }
    })

}

func TestRR_Renew(t *testing.T) {
    bes := lb.BEServers{
        "Address A": new(lb.BEServer),
        "Address B": new(lb.BEServer),
        "Address C": new(lb.BEServer),
        "Address D": new(lb.BEServer),
    }

    testCases := []struct {
        newBes   lb.BEServers
        expected []string
    }{
        {
            newBes: lb.BEServers{
                "Address B": new(lb.BEServer),
                "Address C": new(lb.BEServer),
                "Address D": new(lb.BEServer),
                "Address A": new(lb.BEServer),
                "Address E": new(lb.BEServer), // Add server E.
            },
            expected: []string{"Address A", "Address B", "Address C", "Address D", "Address E"},
        },
        {
            newBes: lb.BEServers{
                "Address B": new(lb.BEServer),
                "Address C": new(lb.BEServer),
                "Address D": new(lb.BEServer), // Delete server A.
            },
            expected: []string{"Address B", "Address C", "Address D"},
        },
        {
            newBes: lb.BEServers{
                "Address B": new(lb.BEServer),
                "Address E": new(lb.BEServer), // Delete server A, C, D and add server E.
            },
            expected: []string{"Address B", "Address E"},
        },
    }

    for _, tc := range testCases {
        rr := NewRR(bes)
        rr.Renew(tc.newBes)

        if !assertEqualSlice(rr.servers, tc.expected) {
            t.Errorf("error renewing server: expected %#v, got %#v.\n", tc.expected, rr.servers)
        }
    }
}

func assertEqualSlice(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for n, ele := range a {
        if ele != b[n] {
            return false
        }
    }

    return true
}
