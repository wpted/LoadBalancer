package lbalgo

import (
    "LoadBalancer/internal/model"
    "errors"
    "net/http"
    "testing"
)

func TestPTC_ChooseServer(t *testing.T) {
    emptyReq := new(http.Request)
    t.Run("No servers", func(t *testing.T) {
        bes := new(model.BEServers)
        rr := NewPTC(bes)
        _, err := rr.ChooseServer(emptyReq)
        if !errors.Is(err, ErrNoServer) {
            t.Errorf("error incorrect error: expected %#v, got %#v.\n", ErrNoServer, err)
        }
    })

    t.Run("One server", func(t *testing.T) {
        bes := &model.BEServers{
            "Address A": &model.BEServer{
                Address:     "Address A",
                Connections: 16,
            },
        }
        ptc := NewPTC(bes)
        addr, err := ptc.ChooseServer(emptyReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }
        if addr != (*bes)["Address A"].Address {
            t.Errorf("error choosing server: expected %#v, got %#v.\n", (*bes)["Address A"].Address, addr)
        }
    })

    t.Run("Correct Backend Server", func(t *testing.T) {
        bes := &model.BEServers{
            "Address A": &model.BEServer{
                Address:     "Address A",
                Connections: 16,
            },
            "Address B": &model.BEServer{
                Address:     "Address B",
                Connections: 22,
            },
            "Address C": &model.BEServer{
                Address:     "Address C",
                Connections: 31,
            },
            "Address D": &model.BEServer{
                Address:     "Address D",
                Connections: 14,
            },
        }

        ptc := NewPTC(bes)
        _, err := ptc.ChooseServer(emptyReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }
    })
}

func TestPTC_Renew(t *testing.T) {
    bes := &model.BEServers{
        "Address A": &model.BEServer{
            Address:     "Address A",
            Connections: 16,
        },
        "Address B": &model.BEServer{
            Address:     "Address B",
            Connections: 22,
        },
        "Address C": &model.BEServer{
            Address:     "Address C",
            Connections: 31,
        },
        "Address D": &model.BEServer{
            Address:     "Address D",
            Connections: 14,
        },
    }

    testCases := []struct {
        newBes   model.BEServers
        expected []string
    }{
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{
                    Address:     "Address B",
                    Connections: 16,
                },
                "Address C": &model.BEServer{
                    Address:     "Address C",
                    Connections: 6,
                },
                "Address D": &model.BEServer{
                    Address:     "Address D",
                    Connections: 19,
                },
                "Address A": &model.BEServer{
                    Address:     "Address A",
                    Connections: 96,
                },
                "Address E": &model.BEServer{
                    Address:     "Address E",
                    Connections: 85,
                }, // Add server E.
            },
            expected: []string{"Address A", "Address B", "Address C", "Address D", "Address E"},
        },
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{
                    Address:     "Address B",
                    Connections: 16,
                },
                "Address C": &model.BEServer{
                    Address:     "Address C",
                    Connections: 6,
                },
                "Address D": &model.BEServer{
                    Address:     "Address D",
                    Connections: 19,
                }, // Delete server A.
            },
            expected: []string{"Address B", "Address C", "Address D"},
        },
        {
            newBes: model.BEServers{
                "Address B": &model.BEServer{
                    Address:     "Address B",
                    Connections: 16,
                },
                "Address E": &model.BEServer{
                    Address:     "Address E",
                    Connections: 85,
                }, // Delete server A, C, D and add server E.
            },
            expected: []string{"Address B", "Address E"},
        },
    }

    for _, tc := range testCases {
        ptc := NewPTC(bes)
        ptc.Renew(tc.newBes)

        addresses := make([]string, 0)

        for _, srv := range ptc.servers {
            addresses = append(addresses, srv.Address)
        }
        if !assertSameElement(addresses, tc.expected) {
            t.Errorf("error renewing server: expected %#v, got %#v.\n", tc.expected, addresses)
        }
    }
}

// assertSameElement checks whether the two string slice contains the same elements.
func assertSameElement(setA, setB []string) bool {
    if len(setA) != len(setB) {
        return false
    }

    dict := make(map[string]struct{})
    for _, address := range setA {
        dict[address] = struct{}{}
    }

    for _, address := range setB {
        if _, ok := dict[address]; !ok {
            return false
        }
    }

    return true
}
