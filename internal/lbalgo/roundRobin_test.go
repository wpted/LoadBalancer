package lbalgo

import (
    "LoadBalancer/internal/model"
    "errors"
    "net/http"
    "testing"
)

func TestRR_ChooseServer(t *testing.T) {

    emptyReq := new(http.Request)
    t.Run("No servers", func(t *testing.T) {
        bes := new(model.BEServers)
        rr := NewRR(bes)
        _, err := rr.ChooseServer(emptyReq)
        if !errors.Is(err, ErrNoServer) {
            t.Errorf("error incorrect error: expected %#v, got %#v.\n", ErrNoServer, err)
        }
    })

    t.Run("Correct Backend Server", func(t *testing.T) {
        bes := &model.BEServers{
            "Address A": new(model.BEServer),
            "Address B": new(model.BEServer),
            "Address C": new(model.BEServer),
            "Address D": new(model.BEServer),
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
    bes := &model.BEServers{
        "Address A": new(model.BEServer),
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer),
        "Address D": new(model.BEServer),
    }

    testCases := []struct {
        newBes   model.BEServers
        expected []string
    }{
        {
            newBes: model.BEServers{
                "Address B": new(model.BEServer),
                "Address C": new(model.BEServer),
                "Address D": new(model.BEServer),
                "Address A": new(model.BEServer),
                "Address E": new(model.BEServer), // Add server E.
            },
            expected: []string{"Address A", "Address B", "Address C", "Address D", "Address E"},
        },
        {
            newBes: model.BEServers{
                "Address B": new(model.BEServer),
                "Address C": new(model.BEServer),
                "Address D": new(model.BEServer), // Delete server A.
            },
            expected: []string{"Address B", "Address C", "Address D"},
        },
        {
            newBes: model.BEServers{
                "Address B": new(model.BEServer),
                "Address E": new(model.BEServer), // Delete server A, C, D and add server E.
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
