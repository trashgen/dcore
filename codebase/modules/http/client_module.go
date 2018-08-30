package http

import (
    "fmt"
    "log"
    "time"
    "net/http"
    "io/ioutil"
    dcconf "dcore/codebase/modules/config"
)

type ClientModule struct {
    client *http.Client
    config *dcconf.ClientConfig
}

func NewClientModule(config *dcconf.ClientConfig) *ClientModule {
    stdHTTPClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}
    return &ClientModule{client:stdHTTPClient, config:config}
}

func (this *ClientModule) RequestReg() string {
    url := fmt.Sprintf("%s/%s", this.config.EntryPoints[0], this.config.Reg.Name)
    response, err := this.getRawContent(url)
    if err != nil {
        err := fmt.Sprintf("Error by getting response 'Look' from Point [%s]: [%s]\n", url, err.Error())
        log.Print(err)
        
        return err
    }

    return fmt.Sprintf("%s\n", response)
}

func (this *ClientModule) RequestLook(maxPoints int, count int) string {
    urls := make([]string, 0, maxPoints)
    for i := 0; i < maxPoints; i++ {
        if count == 0 {
            urls = append(urls, fmt.Sprintf("%s/%s", this.config.EntryPoints[i], this.config.Look.Name))
        } else {
            urls = append(urls, fmt.Sprintf("%s/%s?%s=%d", this.config.EntryPoints[i], this.config.Look.Name, this.config.Look.Param, count))
        }
    }

    var out string
    for _, url := range urls {
        response, err := this.getRawContent(url)
        if err != nil {
            err := fmt.Sprintf("Error by getting response 'Look' from Point [%s]: [%s]\n", url, err.Error())
            log.Print(err)

            return err
        }

        out += fmt.Sprintf("%s\n", response)
    }

    return out
}

func (this *ClientModule) getRawContent(url string, ) (string, error) {
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
