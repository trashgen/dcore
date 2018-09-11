package client

import (
    "log"
    "time"
    "strconv"
    "strings"
    "net/http"
    "io/ioutil"
    dcutil "dcore/codebase/util"
    "dcore/codebase/module/http/command"
    dcconf "dcore/codebase/module/config"
)

type HTTPClient struct {
    httpClient   *http.Client
    clientConfig *dcconf.ClientConfig
}

func NewClientModule(config *dcconf.ClientConfig) *HTTPClient {
    stdHTTPClient := &http.Client{
        Timeout   : time.Second * 11,
        Transport : &http.Transport {
            DisableKeepAlives   : true,
            DisableCompression  : false,
            TLSHandshakeTimeout : time.Second * 11}}
    return &HTTPClient{httpClient: stdHTTPClient, clientConfig: config}
}

func (this *HTTPClient) SendRequestBan(thisKey string, thoseKey string) string {
    return this.getRawContent(command.NewBanRequest(thisKey, thoseKey).RequestToURL(this.clientConfig.EntryPoints[0]))
}

func (this *HTTPClient) SendRequestReg(port int) string {
    return this.getRawContent(command.NewRegRequest(port).RequestToURL(this.clientConfig.EntryPoints[0]))
}

func (this *HTTPClient) RequestLook(key string, numPoints int, numNodes int) map[string]string {
    // TODO : На данный момент нет мультипоинта - посему один запрос на нулевой поинт
    out := make(map[string]string)
    response := this.getRawContent(command.NewLookRequest(key, numPoints, numNodes).RequestToURL(this.clientConfig.EntryPoints[0]))
    triplets := dcutil.ScanString(response, '\t')
    for _, triplet := range triplets {
        values := strings.Split(triplet, "-")
        if values[0] != key {
            out[values[0]] = values[1]
        }
    }
    return out
}

func (this *HTTPClient) RequestPoints(key string, count int) []string {
    // TODO : На данный момент нет мультипоинта - посему один запрос на нулевой поинт
    response := this.getRawContent(command.NewPointsRequest(key, count).RequestToURL(this.clientConfig.EntryPoints[0]))
    return dcutil.ScanString(response, '\t')
}

func (this *HTTPClient) RequestCheck(key string, target string) bool {
    resultValue := this.getRawContent(command.NewCheckRequest(key, target).RequestToURL(this.clientConfig.EntryPoints[0]))
    result, _ := strconv.ParseBool(resultValue)
    return result
}

func (this *HTTPClient) RequestRemove(key string, target string) string {
    return this.getRawContent(command.NewRemoveRequest(key, target).RequestToURL(this.clientConfig.EntryPoints[0]))
}

func (this *HTTPClient) getRawContent(url string) string {
    resp, err := this.httpClient.Get(url)
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
    return string(bodyBytes)
}