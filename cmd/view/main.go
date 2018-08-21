package main

import (
    "log"
    "time"
    "net/http"
    "io/ioutil"
    dcconf "dcore/codebase/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    viewClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}

    for _, request := range config.BuildListAllGetRequestList() {
        if body, err := doListAll(viewClient, request); err == nil {
            log.Printf("Request from SignalServer:\n[%s]\n", body)
        }
    }
}

func doListAll(client *http.Client, url string) (string, error) {
    resp, err := client.Get(url)
    if resp != nil {
        defer resp.Body.Close()
    }
    if err != nil {
        log.Print(err)
        return "", err
    }

    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    return string(bodyBytes), nil
}