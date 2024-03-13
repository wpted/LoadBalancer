package lbalgo

import (
    "LoadBalancer/internal/model"
    "errors"
    "net/http"
    "strings"
)

const (
    LeastConnection    = "LC"
    RoundRobin         = "RR"
    StickyRoundRobin   = "SRR"
    WeightedRoundRobin = "WRR"
    SourceIPHashing    = "SIH"
    PowerOfTwoChoices  = "PTC"
)

var (
    ErrNoServer    = errors.New("error no available server")
    ErrUnknownAlgo = errors.New("error unknown algorithm")
)

type LBAlgo interface {
    ChooseServer(req *http.Request) (string, error)
    Renew(servers model.BEServers)
}

func ChooseAlgo(algoBrief string) (LBAlgo, error) {
    switch strings.ToUpper(algoBrief) {
    case LeastConnection:
        return NewLC(nil), nil
    case RoundRobin:
        return NewRR(nil), nil
    case StickyRoundRobin:
        return NewSRR(nil), nil
    case WeightedRoundRobin:
        return NewWRR(nil), nil
    case SourceIPHashing:
        return NewSIH(nil), nil
    case PowerOfTwoChoices:
        return NewPTC(nil), nil
    default:
        return nil, ErrUnknownAlgo
    }
}

// getClientIP gets the IP of the client. If the client is hided behind proxies or load balancers,
// we get the IP from retrieving the value from X-Forwarded-For header.
// This method is only a demo and shouldn't be used in any production code. Header can be changed after a request is sent.
// The best way is to store all occurring ip for further analytics.
func getClientIP(req *http.Request) string {
    clientIP := req.Header.Get("X-Forwarded-For")
    if clientIP == "" {
        return req.RemoteAddr
    }

    return clientIP
}
