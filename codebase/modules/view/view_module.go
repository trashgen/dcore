package view

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
    dcmod "dcore/codebase/modules/misc"
    dcmisc "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
)

type ResponseMapper func(body string) interface{}

type ViewModule struct {
    Nodes         map[string]*dcmod.NodeID
    client        *http.Client
    config        *dcconf.TotalConfig
    mapCheckRaw   ResponseMapper
    mapRemoveRaw  ResponseMapper
    mapListallRaw ResponseMapper
}

func NewViewModule(config *dcconf.TotalConfig) *ViewModule {
    stdHTTPClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}
    return &ViewModule{
        Nodes         : make(map[string]*dcmod.NodeID),
        client        : stdHTTPClient,
        config        : config,
        mapCheckRaw   : checkResponseHandler,
        mapRemoveRaw  : removeResponseHandler,
        mapListallRaw : listallResponseHandler}
}

func (this *ViewModule) GetRawContent(method string, param string) string {
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

func (this *ViewModule) MapListall(body string) *dcmisc.RequestListall {
    rawResult := this.mapListallRaw(body)
    out, ok := rawResult.(*dcmisc.RequestListall)
    if ! ok {
        log.Fatal("mapListall: inner error")
    }

    return out
}

func (this *ViewModule) MapRemove(body string) *dcmisc.RequestRemove {
    rawResult := this.mapRemoveRaw(body)
    out, ok := rawResult.(*dcmisc.RequestRemove)
    if ! ok {
        log.Fatal("mapRemove: inner error")
    }

    return out
}

func (this *ViewModule) MapCheck(body string) *dcmisc.RequestCheck {
    rawResult := this.mapCheckRaw(body)
    out, ok := rawResult.(*dcmisc.RequestCheck)
    if ! ok {
        log.Fatal("mapCheck: inner error")
    }
    
    return out
}

func (this *ViewModule) buildURL(method string, param string) (string, error) {
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

func (this *ViewModule) getResponseBody(url string) (string, error) {
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

func (this *ViewModule) buildListallURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.ListAll, "count", param)
}

func (this *ViewModule) buildRemoveURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.Remove, "key", param)
}

func (this *ViewModule) buildCheckURL(param string) string {
    return fmt.Sprintf("http://127.0.0.1:30001/%s?%s=%s", this.config.SSCommand.Check, "key", param)
}

func listallResponseHandler(body string) interface{} {
    params := strings.Split(body, "\n")
    out := &dcmisc.RequestListall {
        RequestorID : params[0],
        Nodes       : make(map[string]*dcmisc.NodeID)}
    scanner := bufio.NewScanner(strings.NewReader(params[1]))
    scanner.Split(dcmisc.SplitTabs)
    for scanner.Scan() {
        nodeIDRawData := strings.Split(scanner.Text(), ":")
        port, err := strconv.Atoi(nodeIDRawData[2])
        if err != nil {
            log.Fatalf("mapListallBody: %s\n", err.Error())
        }
        
        out.Nodes[nodeIDRawData[0]] = dcmisc.NewNodeID(nodeIDRawData[0], nodeIDRawData[1], port)
    }
    return out
}

func removeResponseHandler(body string) interface{} {
    param := strings.TrimSuffix(body, "\n")
    opResult, err := strconv.ParseBool(param)
    if err != nil {
        log.Fatal("removeResponseHandler: inner error\n")
    }

    return &dcmisc.RequestRemove{OpResult:opResult}
}

func checkResponseHandler(body string) interface{} {
    param := strings.TrimSuffix(body, "\n")
    opResult, err := strconv.ParseBool(param)
    if err != nil {
        log.Fatal("checkResponseHandler: inner error\n")
    }
    
    return &dcmisc.RequestCheck{OpResult:opResult}
}