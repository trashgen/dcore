package http

import (
    "fmt"
    "log"
    "time"
    "strings"
    "net/http"
    "crypto/md5"
    "encoding/hex"
    dcutil "dcore/codebase/util"
    dcmisc "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
)

type Point struct {
    id        string
    nodes     map[string]*dcutil.ConnectionID
    config    *dcconf.PointConfig
    points    map[string]*dcutil.ConnectionID // not used now
    cmdConfig *dcconf.HTTPCommands
}

func NewPoint(config *dcconf.PointConfig, cmdConfig *dcconf.HTTPCommands) *Point {
    out := &Point{
        nodes     : make(map[string]*dcutil.ConnectionID, 16),
        points    : make(map[string]*dcutil.ConnectionID, 16),
        config    : config,
        cmdConfig : cmdConfig}

    return out
}

func (this *Point) Start() {
    if err := http.ListenAndServe(this.config.FormattedListenPort(), this); err != nil {
        log.Fatalf("Error starting Point: %s", err.Error())
    }
}

func (this *Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Printf("Got request [%s] - [%s]\n", r.URL.Path, r.URL.RawQuery)
    switch r.URL.Path[1:] {
        case this.cmdConfig.Reg.Name:
            this.responseToReg(w, r.RemoteAddr, r.URL.RawQuery)
        case this.cmdConfig.Look.Name:
            this.responseToLook(w, r.URL.RawQuery)
        case this.cmdConfig.Root.Name:
            this.responseToRoot(w, r.URL.RawQuery)
        case this.cmdConfig.Check.Name:
            this.responseToCheck(w, r.URL.RawQuery)
        case this.cmdConfig.Points.Name:
            this.responseToPoints(w, r.URL.RawQuery)
        case this.cmdConfig.Remove.Name:
            this.responseToRemove(w, r.URL.RawQuery)
        default:
            log.Printf("ServeHTTP - Bad HTTP Command: [%s]", r.URL.Path[1:])
    }
}

func (this *Point) responseToReg(w http.ResponseWriter, remoteAddr string, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)

    // Nice form just to test
    if msg := "Reg requires query param 'address'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.Write([]byte(msg))

        return
    }

    request, err := dcmisc.SplitRequestReg(this.cmdConfig.Reg.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))
        
        return
    }

    key := this.calcMD5Key(remoteAddr)
    this.nodes[key] = dcutil.NewConnectionID(key, request.Address)

    w.Write([]byte(fmt.Sprintf("%s", key)))
}

func (this *Point) responseToLook(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    
    request, err := dcmisc.SplitRequestLook(this.cmdConfig.Look.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))

        return
    }

    realCount := request.Count
    if realCount > len(this.nodes) {
        realCount = len(this.nodes)
    }

    sb := strings.Builder{}
    i := 0
    for _, v := range this.nodes {
        sb.WriteString(fmt.Sprintf("%s:%s\t", v.Key, v.Address))
        if i++; i == realCount {
            break
        }
    }

    w.Write([]byte(strings.TrimSuffix(sb.String(), "\t")))
}

func (this *Point) responseToRoot(w http.ResponseWriter, queryParams string) {
    if len(queryParams) > 0 {
        log.Printf("Root: ignore query params: [%s]\n", queryParams)
    }

    sb := strings.Builder{}
    sb.WriteString("<h1>Point help:</h1>")
    sb.WriteString("<b>Root page</b>: You are here now<br>")
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to register on Point. No query params. Key (string) as result<br>", this.cmdConfig.Reg.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Nodes. If query param 'count' here with (int) > 0 as value - then limit number of Nodes to send in Response<br>", this.cmdConfig.Look.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to check if Node is registered at this Point. Key (string) as query param required<br>", this.cmdConfig.Check.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Points. If query param 'count' here with (int) > 0 as value - then limit number of Points to send in Response<br>", this.cmdConfig.Points.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to remove Node. Key (string) as query param required<br>", this.cmdConfig.Remove.Name))

    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(sb.String()))
}

func (this *Point) responseToCheck(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)

    if msg := "Check requires query param 'key'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.Write([]byte(msg))
    
        return
    }

    request, err := dcmisc.SplitRequestCheck(this.cmdConfig.Check.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))
        
        return
    }
    
    _, ok := this.nodes[request.Key]
    if ok {
        w.Write([]byte(fmt.Sprintf("true")))
    } else {
        w.Write([]byte(fmt.Sprintf("false")))
    }
}

func (this *Point) responseToPoints(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)

    request, err := dcmisc.SplitRequestPoints(this.cmdConfig.Points.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))
        
        return
    }
    
    realCount := request.Count
    if realCount > len(this.points) {
        realCount = len(this.points)
    }
    
    sb := strings.Builder{}
    i := 0
    for _, v := range this.points {
        sb.WriteString(fmt.Sprintf("%s:%s\t", v.Key, v.Address))
        if i++; i == realCount {
            break
        }
    }

    w.Write([]byte(strings.TrimSuffix(sb.String(), "\t")))
}

func (this *Point) responseToRemove(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)

    if msg := "Remove requires query param 'key'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.Write([]byte(msg))

        return
    }
    
    request, err := dcmisc.SplitRequestRemove(this.cmdConfig.Remove.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))
    
        return
    }
    
    _, ok := this.nodes[request.Key]
    if ok {
        delete(this.nodes, request.Key)
        w.Write([]byte(fmt.Sprintf("true")))
    } else {
        w.Write([]byte(fmt.Sprintf("false")))
    }

}

func (this *Point) calcMD5Key(addr string) string {
    hash := md5.New()
    hash.Write([]byte(fmt.Sprintf("%s%s%s", addr, time.Now().String(), this.config.SecretPhrase)))
    return hex.EncodeToString(hash.Sum(nil))
}