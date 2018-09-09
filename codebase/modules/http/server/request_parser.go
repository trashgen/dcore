package server

import (
    "log"
    "net/url"
    "dcore/codebase/modules/config"
    dchttp "dcore/codebase/modules/http"
)

// TODO : После создания используется ТОЛЬКО один раз. Пока так. А по факту это Stateless Factory per Point.
type RequestParser struct {
    cmd         string
    cmdConfig   *config.HTTPCommands
    queryParams url.Values
}

func NewRequestParser(cmd string, queryParams url.Values, cmdConfig *config.HTTPCommands) *RequestParser {
    return &RequestParser{cmd: cmd, queryParams: queryParams, cmdConfig: cmdConfig}
}

func (this RequestParser) Ban() *dchttp.RequestBan {
    return this.parse().(*dchttp.RequestBan)
}

func (this RequestParser) Reg() *dchttp.RequestReg {
    return this.parse().(*dchttp.RequestReg)
}

func (this RequestParser) Look() *dchttp.RequestLook {
    return this.parse().(*dchttp.RequestLook)
}

func (this RequestParser) Check() *dchttp.RequestCheck {
    return this.parse().(*dchttp.RequestCheck)
}

func (this RequestParser) Points() *dchttp.RequestPoints {
    return this.parse().(*dchttp.RequestPoints)
}

func (this RequestParser) Remove() *dchttp.RequestRemove {
    return this.parse().(*dchttp.RequestRemove)
}

func (this *RequestParser) parse() interface{} {
    switch this.cmd {
        case this.cmdConfig.Reg.Name:
            return dchttp.CreateRegRequest(this.queryParams)
        case this.cmdConfig.Ban.Name:
            return dchttp.CreateBanRequest(this.queryParams)
        case this.cmdConfig.Look.Name:
            return dchttp.CreateLookRequest(this.queryParams)
        case this.cmdConfig.Check.Name:
            return dchttp.CreateCheckRequest(this.queryParams)
        case this.cmdConfig.Points.Name:
            return dchttp.CreatePointsRequest(this.queryParams)
        case this.cmdConfig.Remove.Name:
            return dchttp.CreateRemoveRequest(this.queryParams)
    }
    log.Fatalf("ServeHTTP - Bad HTTP Command: [%s]", this.cmd)
    return nil
}
