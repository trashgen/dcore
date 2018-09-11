package command

import (
    "fmt"
    "net/url"
    "reflect"
    "dcore/codebase/module/http/persistence"
)

const (
    BanName = "ban"
    ipName  = "ip"
)

type Ban struct {
    Name        string
    ip          *ipDesc
    key         *keyDesc
    redis       *persistence.RedisModule
    paramsCount int
}
func (this Ban) Key() string { return this.key.value }
func (this Ban) IP()  string { return this.ip.value  }

func NewBanRequest(thisKey string, target string) *Ban {
    return &Ban{
        Name        : BanName,
        key:          newKeyDescRequest(thisKey),
        ip:           newIPDescRequest(target),
        paramsCount : 2}
}

func NewBanResponse(redis *persistence.RedisModule) *Ban {
    return &Ban{
        Name        : BanName,
        key         : newKeyDescResponse(),
        redis       : redis,
        ip:           newIPDescResponse(),
        paramsCount : 2}
}

func (this Ban) RequestToURL(pointAddress string) string {
    // ip/ban?key=<string>&ip=<strings>
    return fmt.Sprintf("%s/%s?%s=%s&%s=%s", pointAddress, this.Name,
        this.key.name, this.key.value,
        this.ip.name,  this.ip.value)
}

func (this *Ban) Parse(queryParams url.Values) (err error) {
    if err = checkParamsCount(BanName, queryParams, this.paramsCount); err != nil {
        return err
    }
    if this.key.value, err = tryExtractKeyParam(LookName, keyName, queryParams, this.redis); err != nil {
        return err
    }
    if this.ip.value, err = tryExtractStringParam(LookName, this.ip.name, queryParams); err != nil {
        return err
    }
    return nil
}

func (this Ban) PrepareResponse(params ... interface{}) []byte {
    hddPersist := params[0].(*persistence.MockPersistModule)
    ip := params[1].(string)
    hddPersist.Save(ip)
    return []byte("1")
}



type ipDesc struct {
    name     string
    value    string
    kindType reflect.Kind
}
func newIPDescRequest(target string) *ipDesc {
    return &ipDesc{value: target, name: ipName, kindType: reflect.String}
}
func newIPDescResponse() *ipDesc {
    return &ipDesc{name: ipName, kindType: reflect.String}
}