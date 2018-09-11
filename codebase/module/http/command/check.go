package command

import (
    "fmt"
    "net/url"
    "dcore/codebase/module/http/persistence"
)

const CheckName = "check"

type Check struct {
    Name        string
    key         *keyDesc
    redis       *persistence.RedisModule
    target      *targetDesc
    paramsCount int
}
func (this Check) Key()    string { return this.key.value    }
func (this Check) Target() string { return this.target.value }

func NewCheckRequest(key string, target string) *Check {
    return &Check{
        Name        : CheckName,
        key         : newKeyDescRequest(key),
        target      : newTargetDescRequest(target),
        paramsCount : 2}
}

func NewCheckResponse(redis *persistence.RedisModule) *Check {
    return &Check{
        Name        : CheckName,
        key         : newKeyDescResponse(),
        redis       : redis,
        target      : newTargetDescResponse(),
        paramsCount : 2}
}

func (this Check) RequestToURL(pointAddress string) string {
    // ip/check?key=<string>&target=<strings>
    return fmt.Sprintf("%s/%s?%s=%s&%s=%s", pointAddress, this.Name,
        this.key.name,    this.key.value,
        this.target.name, this.target.value)
}

func (this *Check) Parse(queryParams url.Values) (err error) {
    if err = checkParamsCount(CheckName, queryParams, this.paramsCount); err != nil {
        return err
    }
    if this.key.value, err = tryExtractKeyParam(LookName, keyName, queryParams, this.redis); err != nil {
        return err
    }
    if this.target.value, err = tryExtractKeyParam(LookName, targetName, queryParams, this.redis); err != nil {
        return err
    }
    return nil
}

func (this Check) PrepareResponse(params ... interface{}) []byte {
    target := params[0].(string)
    var result string
    _, ok := this.redis.GetNode(target)
    if ok {
        result = "true"
    } else {
        result = "false"
    }
    return []byte(result)
}