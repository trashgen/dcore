package command

import (
    "fmt"
    "net/url"
    "reflect"
    "strings"
    "dcore/codebase/module/http/persistence"
)

const (
    PointsName = "points"
    countName  = "count"
)

type Points struct {
    Name        string
    key         *keyDesc
    count       *countDesc
    redis       *persistence.RedisModule
    paramsCount int
}
func (this Points) Key()   string { return this.key.value   }
func (this Points) Count() int    { return this.count.value }

func NewPointsRequest(key string, count int) *Points {
    return &Points{
        Name        : LookName,
        key         : newKeyDescRequest(key),
        count       : newCountDescRequest(count),
        paramsCount : 2}
}

func NewPointsResponse(redis *persistence.RedisModule) *Points {
    return &Points{
        Name        : LookName,
        key         : newKeyDescResponse(),
        count       : newCountDescResponse(),
        redis       : redis,
        paramsCount : 2}
}

func (this Points) RequestToURL(pointAddress string) string {
    // ip/points?key=<string>&count=<int>
    return fmt.Sprintf("%s/%s?%s=%s&%s=%d", pointAddress, this.Name,
        this.key.name,   this.key.value,
        this.count.name, this.count.value)
}

func (this *Points) Parse(queryParams url.Values) (err error) {
    if err = checkParamsCount(PointsName, queryParams, this.paramsCount); err != nil {
        return err
    }
    if this.key.value, err = tryExtractKeyParam(LookName, keyName, queryParams, this.redis); err != nil {
        return err
    }
    if this.count.value, err = tryExtractIntParam(LookName, this.count.name, queryParams); err != nil {
        return err
    }
    return nil
}

func (this Points) PrepareResponse(params ... interface{}) []byte {
    nodes := this.redis.GetRandomNodes(this.count.value)
    i := 0
    sb := strings.Builder{}
    for _, v := range nodes {
        sb.WriteString(fmt.Sprintf("%s:%d\t", v.IP, v.Port))
        if i++; i == this.count.value {
            break
        }
    }
    return []byte(sb.String())
}



type countDesc struct {
    name     string
    value    int
    kindType reflect.Kind
}
func newCountDescRequest(count int) *countDesc {
    return &countDesc{value: count, name: countName, kindType: reflect.Int}
}
func newCountDescResponse() *countDesc {
    return &countDesc{name: countName, kindType: reflect.Int}
}
