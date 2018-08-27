package http

import (
    "io"
    "fmt"
    "log"
    "time"
    "errors"
    "strconv"
    "strings"
    "net/http"
    "encoding/hex"
    md52 "crypto/md5"
    dcconf "dcore/codebase/modules/config"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type HTTPModule struct {
    Nodes  map[string]*NodeID
    config *dcconf.TotalConfig
}

func NewHTTPModule(config *dcconf.TotalConfig) *HTTPModule {
    return &HTTPModule{
        Nodes  : make(map[string]*NodeID),
        config : config}
}

func (this *HTTPModule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    cmd := r.URL.Path
    query := r.URL.RawQuery

    handler, err := this.createHandler(cmd[1:], query)
    if err != nil {
        log.Printf(err.Error())
    }

    handler(w, r)
}

func (this *HTTPModule) createHandler(cmd string, query string) (func(w http.ResponseWriter, r *http.Request), error) {
    log.Printf("createHandler: [%s] - [%s]\n", cmd, query)
    switch cmd {
        case this.config.SSCommand.ListAll:
            return handleListall(this, query), nil
        case this.config.SSCommand.Remove:
            return handleRemove(this, query), nil
        case this.config.SSCommand.Check:
            return handleCheck(this, query), nil
    }

    return nil, errors.New(fmt.Sprintf("bad handler command [%s] params [%s]", cmd, query))
}

func handleListall(module *HTTPModule, query string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Connection", "close")
        w.WriteHeader(http.StatusOK)

        sb := strings.Builder{}
        requestorMDKey := module.calcMD5Key(r.RemoteAddr)
        sb.WriteString(fmt.Sprintf("%s\n", requestorMDKey))
        nodeCountToSend := module.getListAllQueryParamValue(query)
        if nodeCountToSend > 0 {
            i := 1
            for _, nodeID := range module.Nodes {
                sb.WriteString(fmt.Sprintf("%s:%s:%d\t", nodeID.ID, nodeID.Address, nodeID.Port))
                if i == len(module.Nodes) {
                    break
                }
        
                i++
            }
        } else {
            for _, nodeID := range module.Nodes {
                sb.WriteString(fmt.Sprintf("%s:%s:%d\t", nodeID.ID, nodeID.Address, nodeID.Port))
            }
        }

        toSend := fmt.Sprintf("%s\n", strings.TrimSuffix(sb.String(), "\t"))
        log.Printf("handleListall response: [%s]", toSend)
        io.WriteString(w, toSend)

        ip, port := module.splitRemoteAddress(r.RemoteAddr)
        module.Nodes[requestorMDKey] = NewNodeID(requestorMDKey, ip, port)
    }
}

func handleRemove(module *HTTPModule, query string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        result := "false\n"
        if _, ok := module.Nodes[module.getKeyQueryParamValue(query)]; ok {
            result = "true\n"
            delete(module.Nodes, module.getKeyQueryParamValue(query))
        }

        w.Header().Set("Connection", "close")
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, result)
    }
}

func handleCheck(module *HTTPModule, query string) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        result := "false\n"
        if _, ok := module.Nodes[module.getKeyQueryParamValue(query)]; ok {
            result = "true\n"
        }

        io.WriteString(w, result)
    }
}

func (this *HTTPModule) splitRemoteAddress(addr string) (string, int) {
    params := strings.Split(addr, ":")
    port, err := strconv.Atoi(params[1])
    if err != nil {
        log.Fatal(err)
    }

    return params[0], port
}

func (this *HTTPModule) calcNumberOfNodesToSend(query string) int {
    out := this.getListAllQueryParamValue(query)
    realNumberOfNodes := len(this.Nodes)
    if out >= realNumberOfNodes {
        out = realNumberOfNodes
    }

    return out
}

func (this *HTTPModule) getListAllQueryParamValue(query string) int {
    params := strings.Split(query, "=")
    count, err := strconv.Atoi(params[1])
    if err != nil {
        log.Fatal(err)
    }

    return count
}

func (this *HTTPModule) getKeyQueryParamValue(query string) string {
    return strings.Split(query, "=")[1]
}

func (this *HTTPModule) calcMD5Key(addr string) string {
    hash := md52.New()
    hash.Write([]byte(fmt.Sprintf("%s%s%s", addr, time.Now().String(), this.config.SecretMD5Phrase)))
    return hex.EncodeToString(hash.Sum(nil))
}