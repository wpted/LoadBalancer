package lbalgo

import (
    "LoadBalancer/internal/model"
    "net/http"
    "testing"
)

func TestSIH_ChooseServer(t *testing.T) {
    bes := &model.BEServers{
        "Address A": new(model.BEServer),
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer),
        "Address D": new(model.BEServer),
    }
    sih := NewSIH(bes)

    testCases := []struct {
        clientReq      *http.Request
        expectedChosen string
    }{
        {clientReq: &http.Request{RemoteAddr: "10.0.0.1"}},
        {clientReq: &http.Request{RemoteAddr: "10.0.0.2"}},
        {clientReq: &http.Request{RemoteAddr: "10.0.0.3"}},
    }

    for _, tc := range testCases {
        chosen, err := sih.ChooseServer(tc.clientReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }

        // Set the first chosen server to test case.
        tc.expectedChosen = chosen

        // Choose again. The two chosen address should be in the same bucket.
        chosen, err = sih.ChooseServer(tc.clientReq)
        if err != nil {
            t.Errorf("error choosing server: got %#v.\n", err)
        }

        expectedBucketNum, _ := sih.exists(tc.expectedChosen)
        bucketNum, _ := sih.exists(chosen)
        if expectedBucketNum != bucketNum {
            t.Errorf("error choosing server: expected bucket num %d, got bucket num %d.\n", expectedBucketNum, bucketNum)
        }
    }
}

func TestSIH_Renew(t *testing.T) {
    bes := &model.BEServers{
        "Address A": new(model.BEServer),
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer),
        "Address D": new(model.BEServer),
    }

    sih := NewSIH(bes)

    newBes := model.BEServers{
        "Address B": new(model.BEServer),
        "Address C": new(model.BEServer), // Delete server A, D.
        "Address E": new(model.BEServer), // Add server E.
    }

    sih.Renew(newBes)

    testCases := []struct {
        address string
        exist   bool
    }{
        {address: "Address A", exist: false},
        {address: "Address B", exist: true},
        {address: "Address C", exist: true},
        {address: "Address D", exist: false},
        {address: "Address E", exist: true},
    }

    for _, tc := range testCases {
        if _, ok := sih.exists(tc.address); ok != tc.exist {
            t.Errorf("error renewing server: expected %t, got %t.\n", tc.exist, ok)
        }
    }
}
