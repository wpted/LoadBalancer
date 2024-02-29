package lbalgo

import (
    "LoadBalancer/internal/model"
    "errors"
    "net/http"
)

const (
    LeastConnection    = "LC"
    RoundRobin         = "RR"
    StickyRoundRobin   = "SRR"
    WeightedRoundRobin = "WRR"
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
    switch algoBrief {
    case LeastConnection:
        return NewLC(nil), nil
    case RoundRobin:
        return NewRR(nil), nil
    case StickyRoundRobin:
        return NewSRR(nil), nil
    case WeightedRoundRobin:
        return NewWRR(nil), nil
    default:
        return nil, ErrUnknownAlgo
    }
}
