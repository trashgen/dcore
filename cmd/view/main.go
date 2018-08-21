package main

import (
    "log"
    "time"
    "net/http"
    "io/ioutil"
    dcconf "dcore/codebase/config"
)

const (
    TimeoutRequest = time.Second * 11
)

func main() {
    config := dcconf.NewConfig()
    config.LoadConfig()

    viewClient := &http.Client{
        Timeout   : TimeoutRequest,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : TimeoutRequest}}
    resp, err := viewClient.Get(buildListAllGetRequest(config))
    if resp != nil {
        defer resp.Body.Close()
    }
    if err != nil {
        log.Fatal(err)
    }
    
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    body := string(bodyBytes)
    log.Printf("Request from SignalServer:\n[%s]\n", body)
}

func buildListAllGetRequest(config *dcconf.Config) string {
    return ""
}