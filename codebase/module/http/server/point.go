package server

import (
    "fmt"
    "log"
    "net/url"
    "net/http"
    "dcore/codebase/module/http/command"
    dcconf "dcore/codebase/module/config"
    "dcore/codebase/module/http/persistence"
)

type Point struct {
    id              string
    redis           *persistence.RedisModule
    config          *dcconf.PointConfig
    hddPersist      persistence.HDDPersist
}

func NewPoint(config *dcconf.PointConfig, hddPersist *persistence.MockPersistModule) *Point {
    return &Point{
        redis      : persistence.NewRedisModule(),
        config     : config,
        hddPersist : hddPersist}
}

func (this *Point) Start() {
    if err := http.ListenAndServe(this.config.FormattedListenPort(), this); err != nil {
        log.Fatalf("Error starting Point: %s", err.Error())
    }
}

func (this *Point) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Connection", "close")
    switch r.URL.Path[1:] {
        case command.RegName:
            this.responseToReg(w, r.URL.Query(), r.RemoteAddr)
        case command.BanName:
            this.responseToBan(w, r.URL.Query())
        case command.LookName:
            this.responseToLook(w, r.URL.Query())
        case command.CheckName:
            this.responseToCheck(w, r.URL.Query())
        case command.PointsName:
            this.responseToPoints(w, r.URL.Query())
        case command.RemoveName:
            this.responseToRemove(w, r.URL.Query())
        default:
            log.Printf("ServeHTTP - Bad HTTP Command: [%s]", r.URL.Path[1:])
    }
}

func (this *Point) sendStatusForbidden(w http.ResponseWriter, url *url.URL) {
    w.WriteHeader(http.StatusForbidden)
    w.Write([]byte(fmt.Sprintf("Access denied for: [%s]\n", url.Hostname())))
}

func (this *Point) responseToBan(w http.ResponseWriter, queryParams url.Values) {
    ban := command.NewBanResponse(this.redis)
    if err := ban.Parse(queryParams); err != nil {
        return
    }
    w.Write(ban.PrepareResponse(this.hddPersist, ban.IP()))
}

func (this *Point) responseToReg(w http.ResponseWriter, queryParams url.Values, remoteAddr string) {
    reg := command.NewRegResponse()
    if err := reg.Parse(queryParams); err != nil {
        return
    }
    w.Write(reg.PrepareResponse(this.redis, this.config.SecretPhrase, remoteAddr))
}

// TODO : Add param Points - use this param only after MultiPoint. Now we just can parse it.
func (this *Point) responseToLook(w http.ResponseWriter, queryParams url.Values) {
    look := command.NewLookResponse(this.redis)
    if err := look.Parse(queryParams); err != nil {
        log.Println(err)
    }
    w.Write(look.PrepareResponse())
}

func (this *Point) responseToPoints(w http.ResponseWriter, queryParams url.Values) {
    points := command.NewLookResponse(this.redis)
    if err := points.Parse(queryParams); err != nil {
        return
    }
    w.Write(points.PrepareResponse())
}

func (this *Point) responseToCheck(w http.ResponseWriter, queryParams url.Values) {
    check := command.NewCheckResponse(this.redis)
    if err := check.Parse(queryParams); err != nil {
        return
    }
    w.Write(check.PrepareResponse(check.Target()))
}

func (this *Point) responseToRemove(w http.ResponseWriter, queryParams url.Values) {
    remove := command.NewRemoveResponse(this.redis)
    if err := remove.Parse(queryParams); err != nil {
        return
    }
    w.Write(remove.PrepareResponse(remove.Target()))
}