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
    id     string
    nodes  map[string]*dcutil.ConnectionID
    config *dcconf.PointConfig
    points map[string]*dcutil.ConnectionID // not used now
}

func NewPoint(config *dcconf.PointConfig) *Point {
    out := &Point{
        nodes  : make(map[string]*dcutil.ConnectionID, 16),
        points : make(map[string]*dcutil.ConnectionID, 16),
        config : config}
    out.nodes["nodeKey1"] = &dcutil.ConnectionID{Key:"nodeKey1", Port:13, IP:"127.0.0.1"}
    out.nodes["nodeKey2"] = &dcutil.ConnectionID{Key:"nodeKey2", Port:13, IP:"127.0.0.1"}
    out.points["pointKey1"] = &dcutil.ConnectionID{Key:"pointKey1", Port:666, IP:"127.0.0.721"}
    out.points["pointKey2"] = &dcutil.ConnectionID{Key:"pointKey2", Port:666, IP:"127.0.0.721"}

    return out
}

func (this *Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.URL.Path[1:] {
        case this.config.Reg.Name:
            this.responseToReg(w, r.RemoteAddr, r.URL.RawQuery)
        case this.config.Look.Name:
            this.responseToLook(w, r.URL.RawQuery)
        case this.config.Root.Name:
            this.responseToRoot(w, r.URL.RawQuery)
        case this.config.Check.Name:
            this.responseToCheck(w, r.URL.RawQuery)
        case this.config.Points.Name:
            this.responseToPoints(w, r.URL.RawQuery)
        case this.config.Remove.Name:
            this.responseToRemove(w, r.URL.RawQuery)
        default:
            log.Printf("ServeHTTP - Bad HTTP Command: [%s]", r.URL.Path[1:])
    }
}

func (this *Point) responseToReg(w http.ResponseWriter, remoteAddr string, queryParams string) {
    if len(queryParams) > 0 {
        log.Printf("Reg: ignore query params: [%s]\n", queryParams)
    }

    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("%s\n", this.calcMD5Key(remoteAddr))))
}

func (this *Point) responseToLook(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    
    request, err := dcmisc.SplitRequestLook(this.config.Look.Param, queryParams)
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
        sb.WriteString(fmt.Sprintf("%s:%s:%d\t", v.Key, v.IP, v.Port))
        if i++; i == realCount {
            break
        }
    }

    toSend := fmt.Sprintf("%s\n", strings.TrimSuffix(sb.String(), "\t"))
    w.Write([]byte(toSend))
}

func (this *Point) responseToRoot(w http.ResponseWriter, queryParams string) {
    if len(queryParams) > 0 {
        log.Printf("Root: ignore query params: [%s]\n", queryParams)
    }

    sb := strings.Builder{}
    sb.WriteString("<h1>Point help:</h1>")
    sb.WriteString("<b>Root page</b>: You are here now<br>")
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to register on Point. No query params. Key (string) as result<br>", this.config.Reg.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Nodes. If query param 'count' here with (int) > 0 as value - then limit number of Nodes to send in Response<br>", this.config.Look.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to check if Node is registered at this Point. Key (string) as query param required<br>", this.config.Check.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Points. If query param 'count' here with (int) > 0 as value - then limit number of Points to send in Response<br>", this.config.Points.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to remove Node. Key (string) as query param required<br>", this.config.Remove.Name))

    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(sb.String()))
}

func (this *Point) responseToCheck(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)

    request, err := dcmisc.SplitRequestCheck(this.config.Check.Param, queryParams)
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

    request, err := dcmisc.SplitRequestPoints(this.config.Points.Param, queryParams)
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
        sb.WriteString(fmt.Sprintf("%s:%s:%d\t", v.Key, v.IP, v.Port))
        if i++; i == realCount {
            break
        }
    }

    toSend := fmt.Sprintf("%s\n", strings.TrimSuffix(sb.String(), "\t"))
    w.Write([]byte(toSend))
}

func (this *Point) responseToRemove(w http.ResponseWriter, queryParams string) {
    request, err := dcmisc.SplitRequestRemove(this.config.Remove.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.Write([]byte(err.Error()))
        
        return
    }
    
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
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