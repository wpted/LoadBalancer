package lbalgo

import (
    "LoadBalancer/internal/lb"
    "errors"
    "testing"
)

func TestRR_ChooseServer(t *testing.T) {

    t.Run("No servers", func(t *testing.T) {
        bes := new(lb.BEServers)
        rr := New(*bes)
        _, err := rr.ChooseServer()
        if !errors.Is(err, ErrNoServer) {
            t.Errorf("error incorrect error: expected %#v, got %#v.\n", ErrNoServer, err)
        }
    })

    t.Run("Correct Backend Server", func(t *testing.T) {
        bes := lb.BEServers{
            "Address A": struct{}{},
            "Address B": struct{}{},
            "Address C": struct{}{},
            "Address D": struct{}{},
        }

        // After 1 rotate.
        expectedChosen := "Address A"
        expectedRR := []string{"Address B", "Address C", "Address D", "Address A"}

        rr := New(bes)
        res, err := rr.ChooseServer()
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
        "Address A": struct{}{},
        "Address B": struct{}{},
        "Address C": struct{}{},
        "Address D": struct{}{},
    }

    testCases := []struct {
        newBes   lb.BEServers
        expected []string
    }{
        {
            newBes: lb.BEServers{
                "Address B": struct{}{},
                "Address C": struct{}{},
                "Address D": struct{}{},
                "Address A": struct{}{},
                "Address E": struct{}{}, // Add server E.
            },
            expected: []string{"Address A", "Address B", "Address C", "Address D", "Address E"},
        },
        {
            newBes: lb.BEServers{
                "Address B": struct{}{},
                "Address C": struct{}{},
                "Address D": struct{}{}, // Delete server A.
            },
            expected: []string{"Address B", "Address C", "Address D"},
        },
        {
            newBes: lb.BEServers{
                "Address B": struct{}{},
                "Address E": struct{}{}, // Delete server A, C, D and add server E.
            },
            expected: []string{"Address B", "Address E"},
        },
    }

    for _, tc := range testCases {
        rr := New(bes)
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
