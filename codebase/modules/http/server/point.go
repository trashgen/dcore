package server

import (
    "log"
    "net/http"
    dcutil "dcore/codebase/util"
    "dcore/codebase/modules/persistance"
    dcconf "dcore/codebase/modules/config"
    dcuhttp "dcore/codebase/util/http/server"
)

type Point struct {
    id              string
    redis           *persistance.RedisModule
    config          *dcconf.PointConfig
    cmdConfig       *dcconf.HTTPCommands
    checkBlackList  chan string
    resultBlackList chan bool
}

func NewPoint(config *dcconf.PointConfig, cmdConfig *dcconf.HTTPCommands) *Point {
    return &Point{
        redis           : persistance.NewRedisModule(),
        config          : config,
        cmdConfig       : cmdConfig,
        checkBlackList  : make(chan string),
        resultBlackList : make(chan bool)}
}

func (this *Point) Start() {
    go func() {
        postgres := persistance.NewBlackListModule()
        defer postgres.Close()
        for ip := range this.checkBlackList {
            ip := ip
            this.resultBlackList <- postgres.CheckInBlackList(ip)
        }
    }()
    if err := http.ListenAndServe(this.config.FormattedListenPort(), this); err != nil {
        log.Fatalf("Error starting Point: %s", err.Error())
    }
}

func (this *Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    this.checkBlackList <- r.URL.Hostname()
    if <-this.resultBlackList {
        // TODO : возможно что-то интересное будет здесь. Но пока этого достаточно - ничего плохиш не получит!
        return
    }
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

    if msg := "Reg requires query param 'port'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.Write([]byte(msg))
        return
    }

    request, err := dcutil.SplitRequestReg(this.cmdConfig.Reg.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    response, key, ip := dcuhttp.BuildRegResponse(remoteAddr, this.config.SecretPhrase)
    this.redis.AddNode(key, ip, request.Port)
    w.Write(response)
}

func (this *Point) responseToLook(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    request, err := dcutil.SplitRequestLook(this.cmdConfig.Look.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.Write(dcuhttp.BuildLookOrPointsResponse(this.redis.GetAllNodes(), request.Count))
}

func (this *Point) responseToRoot(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    w.Write(dcuhttp.BuildRootResponse(this.cmdConfig))
}

func (this *Point) responseToCheck(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    if msg := "Check requires query param 'key'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.Write([]byte(msg))
    
        return
    }

    request, err := dcutil.SplitRequestCheck(this.cmdConfig.Check.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        
        return
    }
    
    response, _ := dcuhttp.BuildCheckOrRemoveResponse(this.redis.GetAllNodes(), request.Key)
    w.Write(response)
}

func (this *Point) responseToPoints(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    request, err := dcutil.SplitRequestPoints(this.cmdConfig.Points.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    w.Write(dcuhttp.BuildLookOrPointsResponse(this.redis.GetAllNodes(), request.Count))
}

func (this *Point) responseToRemove(w http.ResponseWriter, queryParams string) {
    w.Header().Set("Connection", "close")
    if msg := "Remove requires query param 'key'\n"; len(queryParams) == 0 {
        log.Print(msg)
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(msg))

        return
    }

    request, err := dcutil.SplitRequestRemove(this.cmdConfig.Remove.Param, queryParams)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
    
        return
    }

    response, has := dcuhttp.BuildCheckOrRemoveResponse(this.redis.GetAllNodes(), request.Key)
    if has {
        this.redis.RemoveNode(request.Key)
    }

    w.Write(response)
}