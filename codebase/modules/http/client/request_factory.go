package client

import (
    "log"
    "io/ioutil"
    "net/http"
    "dcore/codebase/modules/config"
)

type RequestFactory struct {
    client    *http.Client
    cmdConfig *config.HTTPCommands
}

func NewRequestFactory(client *http.Client, cmdConfig *config.HTTPCommands) *RequestFactory {
    return &RequestFactory{client: client, cmdConfig: cmdConfig}
}

func (this *RequestFactory) CreateReg()

func (this *RequestFactory) getRawContent(url string) (string, error) {
    resp, err := this.client.Get(url)
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
