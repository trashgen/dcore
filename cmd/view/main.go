package main

import (
    "log"
    "time"
    "net/http"
    "io/ioutil"
)

const (
    TimeoutRequest = time.Second * 11
    SignalServerRequest = "http://127.0.0.1:30001/listall"
)

func main() {
    viewClient := &http.Client{
        Timeout   : TimeoutRequest,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : TimeoutRequest}}
    resp, err := viewClient.Get(SignalServerRequest)
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