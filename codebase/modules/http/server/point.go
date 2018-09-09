package server

import (
    "fmt"
    "log"
    "net/url"
    "net/http"
    dcconf "dcore/codebase/modules/config"
)

type Point struct {
    id              string
    redis           *RedisModule
    config          *dcconf.PointConfig
    cmdConfig       *dcconf.HTTPCommands
    hddPersist      HDDPersist
    checkBlackList  chan string
    resultBlackList chan bool
}

func NewPoint(config *dcconf.PointConfig, cmdConfig *dcconf.HTTPCommands, hddPersist HDDPersist) *Point {
    return &Point{
        redis           : NewRedisModule(),
        config          : config,
        cmdConfig       : cmdConfig,
        hddPersist      : hddPersist,
        checkBlackList  : make(chan string),
        resultBlackList : make(chan bool)}
}

func (this *Point) Start() {
    go func() {
        defer this.hddPersist.Close()
        for ip := range this.checkBlackList {
            ip := ip
            this.resultBlackList <- this.hddPersist.CheckExists(ip)
        }
    }()
    if err := http.ListenAndServe(this.config.FormattedListenPort(), this); err != nil {
        log.Fatalf("Error starting Point: %s", err.Error())
    }
}

func (this *Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Connection", "close")
    if this.checkBlackList <- r.URL.Hostname(); <-this.resultBlackList {
        this.sendStatusForbidden(w, r.URL)
        return
    }
    queryParams, ok := this.validateRequest(r.URL)
    if ! ok {
        this.sendStatusBadRequest(w, r.URL)
        return
    }
    // TODO : Печально, что для обработки вышло 2 свича на одно и тоже, но пока пусть будет так.
    switch r.URL.Path[1:] {
        case this.cmdConfig.Reg.Name:
            this.responseToReg(w, queryParams, r.RemoteAddr)
        case this.cmdConfig.Ban.Name:
            this.responseToBan(w, queryParams)
        case this.cmdConfig.Look.Name:
            this.responseToLook(w, queryParams)
        case this.cmdConfig.Root.Name:
            this.responseToRoot(w)
        case this.cmdConfig.Check.Name:
            this.responseToCheck(w, queryParams)
        case this.cmdConfig.Points.Name:
            this.responseToPoints(w, queryParams)
        case this.cmdConfig.Remove.Name:
            this.responseToRemove(w, queryParams)
        default:
            log.Printf("ServeHTTP - Bad HTTP Command: [%s]", r.URL.Path[1:])
    }
}

func (this *Point) sendStatusForbidden(w http.ResponseWriter, url *url.URL) {
    w.WriteHeader(http.StatusForbidden)
    w.Write([]byte(fmt.Sprintf("Access denied for: [%s]\n", url.Hostname())))
}

func (this *Point) sendStatusBadRequest(w http.ResponseWriter, url *url.URL) {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(fmt.Sprintf("Request is not valid format: [%s]?[%s]\n", url.Path[1:], url.RawQuery)))
}

func (this *Point) validateRequest(request *url.URL) (url.Values, bool) {
    command := request.Path[1:]
    queryParams, err := url.ParseQuery(request.RawQuery)
    if err != nil {
        return nil, false
    }
    if ! this.cmdConfig.IsValidRequest(command, queryParams) {
        return nil, false
    }

    return queryParams, true
}

func (this *Point) responseToBan(w http.ResponseWriter, queryParams url.Values) {
    request := NewRequestParser(this.cmdConfig.Ban.Name, queryParams, this.cmdConfig).Ban()
    if _, has := this.redis.GetNode(request.Key); has {
        this.hddPersist.Save(request.Target)
        w.Write(BuildBanResponse())
    }
}

func (this *Point) responseToReg(w http.ResponseWriter, queryParams url.Values, remoteAddr string) {
    request := NewRequestParser(this.cmdConfig.Reg.Name, queryParams, this.cmdConfig).Reg()
    response, key, ip := BuildRegResponse(remoteAddr, this.config.SecretPhrase)
    this.redis.AddNode(key, ip, request.Port)
    w.Write(response)
}

// TODO : Add param Points - use this param only after MultiPoint. Now we just can parse it.
func (this *Point) responseToLook(w http.ResponseWriter, queryParams url.Values) {
    request := NewRequestParser(this.cmdConfig.Look.Name, queryParams, this.cmdConfig).Look()
    w.Write(BuildLookOrPointsResponse(this.redis.GetAllNodes(), request.Nodes))
}

func (this *Point) responseToRoot(w http.ResponseWriter) {
    w.Write(BuildRootResponse(this.cmdConfig))
}

func (this *Point) responseToCheck(w http.ResponseWriter, queryParams url.Values) {
    request := NewRequestParser(this.cmdConfig.Check.Name, queryParams, this.cmdConfig).Check()
    if _, has := this.redis.GetNode(request.Key); has {
        response, _ := BuildCheckOrRemoveResponse(this.redis.GetAllNodes(), request.Target)
        w.Write(response)
    }
}

func (this *Point) responseToPoints(w http.ResponseWriter, queryParams url.Values) {
    request := NewRequestParser(this.cmdConfig.Points.Name, queryParams, this.cmdConfig).Points()
    w.Write(BuildLookOrPointsResponse(this.redis.GetAllNodes(), request.Count))
}

func (this *Point) responseToRemove(w http.ResponseWriter, queryParams url.Values) {
    request := NewRequestParser(this.cmdConfig.Remove.Name, queryParams, this.cmdConfig).Remove()
    if _, has := this.redis.GetNode(request.Key); has {
        response, has := BuildCheckOrRemoveResponse(this.redis.GetAllNodes(), request.Target)
        if has {
            this.redis.RemoveNode(request.Key)
        }
        w.Write(response)
    }
}