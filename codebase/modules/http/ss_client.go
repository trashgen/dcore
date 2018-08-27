package http

import (
    "fmt"
    "log"
    "time"
    "bufio"
    "errors"
    "strconv"
    "strings"
    "net/http"
    "io/ioutil"
    dcmisc "dcore/codebase/modules/util"
    dcconf "dcore/codebase/modules/config"
)

type ResponseMapper func(body string) interface{}

type SSClient struct {
    Nodes         map[string]*NodeID
    client        *http.Client
    config        *dcconf.TotalConfig
    mapCheckRaw   ResponseMapper
    mapRemoveRaw  ResponseMapper
    mapListallRaw ResponseMapper
}

func NewSSClient(config *dcconf.TotalConfig) *SSClient {
    stdHTTPClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}
    return &SSClient{
        Nodes         : make(map[string]*NodeID),
        client        : stdHTTPClient,
        config        : config,
        mapCheckRaw   : checkResponseHandler,
        mapRemoveRaw  : removeResponseHandler,
        mapListallRaw : listallResponseHandler}
}

func (this *SSClient) GetRawContent(method string, param string) string {
    url, err := this.buildURL(method, param)
    if err != nil {
        log.Fatalf("Bad cmd params:\n\tmethod = [%s]\n\tQueryParam = [%s]\n", method, param)
    }

    out, err := this.getResponseBody(url)
    if err != nil {
        log.Fatal(err.Error())
    }

    return out
}

func (this *SSClient) MapListall(body string) *ResponseListall {
    rawResult := this.mapListallRaw(body)
    out, ok := rawResult.(*ResponseListall)
    if ! ok {
        log.Fatal("mapListall: inner error")
    }

    return out
}

func (this *SSClient) MapRemove(body string) *ResponseRemove {
    rawResult := this.mapRemoveRaw(body)
    out, ok := rawResult.(*ResponseRemove)
    if ! ok {
        log.Fatal("mapRemove: inner error")
    }

    return out
}

func (this *SSClient) MapCheck(body string) *ResponseCheck {
    rawResult := this.mapCheckRaw(body)
    out, ok := rawResult.(*ResponseCheck)
    if ! ok {
        log.Fatal("mapCheck: inner error")
    }
    
    return out
}

func (this *SSClient) buildURL(method string, param string) (string, error) {
    switch method {
        case this.config.SSCommand.ListAll:
            return this.buildListallURL(param), nil
        case this.config.SSCommand.Remove:
            return this.buildRemoveURL(param), nil
        case this.config.SSCommand.Check:
            return this.buildCheckURL(param), nil
    }

    return "", errors.New(fmt.Sprintf("bad method [%s", method))
}

func (this *SSClient) getResponseBody(url string) (string, error) {
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

func (this *SSClient) buildListallURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.ListAll, "count", param)
}

func (this *SSClient) buildRemoveURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.Remove, "key", param)
}

func (this *SSClient) buildCheckURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.Check, "key", param)
}

func listallResponseHandler(body string) interface{} {
    // TODO : не добавлять в список VIEW !!!!!!!!!!!
    params := strings.Split(body, "\n")
    out := &ResponseListall{
        RequestorID : params[0],
        Nodes       : make(map[string]*NodeID)}
    scanner := bufio.NewScanner(strings.NewReader(params[1]))
    scanner.Split(dcmisc.SplitTabs)
    for scanner.Scan() {
        nodeIDRawData := strings.Split(scanner.Text(), ":")
        port, err := strconv.Atoi(nodeIDRawData[2])
        if err != nil {
            log.Fatalf("mapListallBody: %s\n", err.Error())
        }
        
        out.Nodes[nodeIDRawData[0]] = NewNodeID(nodeIDRawData[0], nodeIDRawData[1], port)
    }
    return out
}

func removeResponseHandler(body string) interface{} {
    param := strings.TrimSuffix(body, "\n")
    opResult, err := strconv.ParseBool(param)
    if err != nil {
        log.Fatal("removeResponseHandler: inner error\n")
    }

    return &ResponseRemove{OpResult:opResult}
}

func checkResponseHandler(body string) interface{} {
    param := strings.TrimSuffix(body, "\n")
    opResult, err := strconv.ParseBool(param)
    if err != nil {
        log.Fatal("checkResponseHandler: inner error\n")
    }
    
    return &ResponseCheck{OpResult:opResult}
}