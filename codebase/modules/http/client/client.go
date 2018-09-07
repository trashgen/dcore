package client

import (
    "fmt"
    "log"
    "time"
    "errors"
    "strconv"
    "strings"
    "net/http"
    "io/ioutil"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
)

type HTTPClient struct {
    client       *http.Client
    cmdConfig    *dcconf.HTTPCommands
    clientConfig *dcconf.ClientConfig
}

func NewClientModule(config *dcconf.ClientConfig, cmdConfig *dcconf.HTTPCommands) *HTTPClient {
    stdHTTPClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}
    return &HTTPClient{client: stdHTTPClient, clientConfig: config, cmdConfig: cmdConfig}
}

func (this *HTTPClient) RequestReg(port int) (string, error) {
    url := buildURLWithParams(this.clientConfig.EntryPoints[0], &this.cmdConfig.Reg, port)
    response, err := this.getRawContent(url)
    if err != nil {
        return "", errors.New(fmt.Sprintf("Error by getting response 'Reg' from Point [%s]: [%s]\n", url, err.Error()))
    }
    return fmt.Sprintf("%s", response), nil
}

// TODO : Look must have 2 var query params
func (this *HTTPClient) RequestLook(maxPoints int, count int) []string {
    out := make([]string, 0, count)
    urls := make([]string, 0, maxPoints)
    for i := 0; i < maxPoints; i++ {
        if count == 0 {
            urls = append(urls, buildURLNoParams(this.clientConfig.EntryPoints[i], &this.cmdConfig.Look))
        } else {
            urls = append(urls, buildURLWithParams(this.clientConfig.EntryPoints[i], &this.cmdConfig.Look, count))
        }
    }
    for _, url := range urls {
        response, err := this.getRawContent(url)
        if err != nil {
            log.Printf("Error by getting response 'Look' from Point [%s]: [%s]\n", url, err.Error())
            continue
        }
        regHosts := dcutil.ScanString(strings.TrimSuffix(response, "\n"), '\t')
        out = append(out, regHosts...)
    }
    return out
}

// TODO : Points must have 2 var query params
func (this *HTTPClient) RequestPoints(maxPoints int, count int) string {
    urls := make([]string, 0, maxPoints)
    for i := 0; i < maxPoints; i++ {
        if count == 0 {
            urls = append(urls, buildURLNoParams(this.clientConfig.EntryPoints[i], &this.cmdConfig.Points))
        } else {
            urls = append(urls, buildURLWithParams(this.clientConfig.EntryPoints[i], &this.cmdConfig.Points, count))
        }
    }
    var out string
    for _, url := range urls {
        response, err := this.getRawContent(url)
        if err != nil {
            errDesc := fmt.Sprintf("Error by getting response 'Points' from Point [%s]: [%s]\n", url, err.Error())
            log.Print(errDesc)
            
            return errDesc
        }
        out += fmt.Sprintf("%s\n", response)
    }
    return out
}

func (this *HTTPClient) RequestRoot() string {
    url := buildURLNoParams(this.clientConfig.EntryPoints[0], &this.cmdConfig.Root)
    response, err := this.getRawContent(url)
    if err != nil {
        err := fmt.Sprintf("Error by getting response 'Root' from Point [%s]: [%s]\n", url, err.Error())
        log.Print(err)
        
        return err
    }
    return fmt.Sprintf("%s\n", response)
}

func (this *HTTPClient) RequestCheck(key string) bool {
    url := buildURLWithParams(this.clientConfig.EntryPoints[0], &this.cmdConfig.Check, key)
    response, err := this.getRawContent(url)
    if err != nil {
        log.Print(fmt.Sprintf("Error by getting response 'Check' from Point [%s]: [%s]\n", url, err.Error()))
        return false
    }
    status, err := strconv.ParseBool(response)
    if err != nil {
        log.Fatalln(err.Error())
    }
    return status
}

func (this *HTTPClient) RequestRemove(key string) string {
    url := buildURLWithParams(this.clientConfig.EntryPoints[0], &this.cmdConfig.Remove, key)
    response, err := this.getRawContent(url)
    if err != nil {
        err := fmt.Sprintf("Error by getting response 'Remove' from Point [%s]: [%s]\n", url, err.Error())
        log.Print(err)
        
        return err
    }
    return fmt.Sprintf("%s\n", response)
}

func (this *HTTPClient) getRawContent(url string) (string, error) {
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

func buildURLNoParams(pointAddress string, commandDesc *dcconf.CommandDesc) string {
    return fmt.Sprintf("%s/%s", pointAddress, commandDesc.Name)
}

func buildURLWithParams(pointAddress string, commandDesc *dcconf.CommandDesc, param interface{}) (out string) {
    maybeString, ok := param.(string)
    if ok {
        return fmt.Sprintf("%s/%s?%s=%s", pointAddress, commandDesc.Name, commandDesc.Param, maybeString)
    }
    maybeInt, ok := param.(int)
    if ok {
        return fmt.Sprintf("%s/%s?%s=%d", pointAddress, commandDesc.Name, commandDesc.Param, maybeInt)
    }
    log.Fatalf("URL with bad param [%#v]\n", param)
    return out
}