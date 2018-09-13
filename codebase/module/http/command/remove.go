package command

import (
    "fmt"
    "net/url"
    "dcore/codebase/module/http/persistence"
)

const RemoveName = "remove"

type Remove struct {
    Name        string
    key         *keyDesc
    redis       *persistence.RedisModule
    target      *targetDesc
    paramsCount int
}
func (this Remove) Key()    string { return this.key.value    }
func (this Remove) Target() string { return this.target.value }

func NewRemoveRequest(key string, target string) *Remove {
    return &Remove{
        Name        : RemoveName,
        key         : newKeyDescRequest(key),
        target      : newTargetDescRequest(target),
        paramsCount : 2}
}

func NewRemoveResponse(redis *persistence.RedisModule) *Remove {
    return &Remove{
        Name        : RemoveName,
        key         : newKeyDescResponse(),
        redis       : redis,
        target      : newTargetDescResponse(),
        paramsCount : 2}
}

func (this Remove) RequestToURL(pointAddress string) string {
    // ip/remove?key=<string>&target=<strings>
    return fmt.Sprintf("%s/%s?%s=%s&%s=%s", pointAddress, this.Name,
        this.key.name,    this.key.value,
        this.target.name, this.target.value)
}

func (this *Remove) Parse(queryParams url.Values) (err error) {
    if err = checkParamsCount(RemoveName, queryParams, this.paramsCount); err != nil {
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

func (this Remove) PrepareResponse(params ... interface{}) []byte {
    target := params[0].(string)
    var result string
    if this.redis.NodeExists(target) {
        result = "true"
        this.redis.RemoveNode(target)
    } else {
        result = "false"
    }
    return []byte(result)
}