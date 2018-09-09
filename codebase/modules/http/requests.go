package http

import (
    "net/url"
    "strconv"
    "dcore/codebase/modules/config"
    "fmt"
)

type RequestBan struct {Key string; Target string}
func CreateBanRequest(queryParams url.Values) *RequestBan {
    return &RequestBan{
        Key    : queryParams.Get(config.HTTPParamKey),
        Target : queryParams.Get(config.HTTPParamTarget)}
}



type RequestReg struct {Port int; address string; commandDesc *config.CommandDesc}
func NewRequestReg(address string, port int, commandDesc *config.CommandDesc) *RequestReg {
    return &RequestReg{address: address, Port: port, commandDesc: commandDesc}
}
func CreateRegRequest(queryParams url.Values) *RequestReg {
    val, _ := strconv.Atoi(queryParams.Get(config.HTTPParamPort))
    return &RequestReg{Port: val}
}
func (this *RequestReg) String() string {
    return fmt.Sprintf("%s/%s?%s=%d", this.address, this.Port)
}



type RequestLook struct {Points int; Nodes int}
func CreateLookRequest(queryParams url.Values) *RequestLook {
    nodes, _  := strconv.Atoi(queryParams.Get(config.HTTPParamNodes))
    points, _ := strconv.Atoi(queryParams.Get(config.HTTPParamPoints))
    return &RequestLook{Points: points, Nodes: nodes}
}

type RequestCheck struct {Key string; Target string}
func CreateCheckRequest(queryParams url.Values) *RequestCheck {
    return &RequestCheck{
        Key    : queryParams.Get(config.HTTPParamKey),
        Target : queryParams.Get(config.HTTPParamTarget)}
}

type RequestPoints struct {Count int}
func CreatePointsRequest(queryParams url.Values) *RequestPoints {
    val, _ := strconv.Atoi(queryParams.Get(config.HTTPParamCount))
    return &RequestPoints{Count: val}
}

type RequestRemove struct {Key string; Target string}
func CreateRemoveRequest(queryParams url.Values) *RequestRemove {
    return &RequestRemove{
        Key    : queryParams.Get(config.HTTPParamKey),
        Target : queryParams.Get(config.HTTPParamTarget)}
}