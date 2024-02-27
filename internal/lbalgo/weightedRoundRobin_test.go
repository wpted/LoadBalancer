package lbalgo

import (
    "LoadBalancer/internal/lb"
    "net/http"
)

type WRR struct {
}

func (w *WRR) ChooseServer(req *http.Request) (string, error) {
    return "", nil
}

func (w *WRR) Renew(servers lb.BEServers) {

}
