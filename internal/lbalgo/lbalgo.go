package lbalgo

import (
    "LoadBalancer/internal/lb"
    "errors"
    "net/http"
)

type LBAlgo interface {
    ChooseServer(req *http.Request) (string, error)
    Renew(servers lb.BEServers)
}

var ErrNoServer = errors.New("error no available server")
