package lbalgo

import "LoadBalancer/internal/lb"

type LBAlgo interface {
    ChooseServer() (string, error)
    Renew(servers lb.BEServers)
}
