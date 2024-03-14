package lbalgo

import (
    "LoadBalancer/internal/model"
    "net/http"
)

type CH struct {
}

func NewCH(backendServers *model.BEServers) *CH {
    return nil
}

func (c *CH) ChooseServer(req *http.Request) (string, error) {
    return "", nil
}

func (c *CH) Renew(currentlyHealthyServers model.BEServer) {

}
