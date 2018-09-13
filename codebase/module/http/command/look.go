package command

import (
    "fmt"
    "net/url"
    "reflect"
    "dcore/codebase/module/http/persistence"
)

const (
    LookName      = "look"
    numNodesName  = "nodes"
    numPointsName = "points"
)

type Look struct {
    Name        string
    key         *keyDesc
    redis       *persistence.RedisModule
    numNodes    *numNodesDesc
    numPoints   *numPointsDesc
    paramsCount int
}
func (this Look) Key()       string { return this.key.value       }
func (this Look) NumNodes()  int    { return this.numNodes.value  }
func (this Look) NumPoints() int    { return this.numPoints.value }

func NewLookRequest(key string, numPoints int, numNodes int) *Look {
    return &Look{
        Name        : LookName,
        key         : newKeyDescRequest(key),
        numNodes    : newNumNodesDescRequest(numNodes),
        numPoints   : newNumPointsDescRequest(numPoints),
        paramsCount : 3}
}

func NewLookResponse(redis *persistence.RedisModule) *Look {
    return &Look{
        Name        : LookName,
        key         : newKeyDescResponse(),
        redis       : redis,
        numNodes    : newNumNodesDescResponse(),
        numPoints   : newNumPointsDescResponse(),
        paramsCount : 3}
}

func (this Look) RequestToURL(pointAddress string) string {
    // ip/look?key=<string>&numpoints=<int>&numnodes=<int>
    return fmt.Sprintf("%s/%s?%s=%s&%s=%d&%s=%d", pointAddress, this.Name,
        this.key.name,       this.key.value,
        this.numPoints.name, this.numPoints.value,
        this.numNodes.name,  this.numNodes.value)
}

func (this *Look) Parse(queryParams url.Values) (err error) {
    if err = checkParamsCount(LookName, queryParams, this.paramsCount); err != nil {
        return err
    }
    if this.key.value, err = tryExtractKeyParam(LookName, keyName, queryParams, this.redis); err != nil {
        return err
    }
    if this.numPoints.value, err = tryExtractIntParam(LookName, this.numPoints.name, queryParams); err != nil {
        return err
    }
    if this.numNodes.value, err = tryExtractIntParam(LookName, this.numNodes.name, queryParams); err != nil {
        return err
    }
    return nil
}

func (this Look) PrepareResponse(params... interface{}) []byte {
    return this.redis.GetRandomNodes(this.numNodes.value)
    //return this.redis.GetAllNodes()
}

///////////////////////////////////////////////////////////////////////////////

type numPointsDesc struct {
    name     string
    value    int
    kindType reflect.Kind
}
func newNumPointsDescRequest(numPoints int) *numPointsDesc {
    return &numPointsDesc{value: numPoints, name: numPointsName, kindType: reflect.Int}
}
func newNumPointsDescResponse() *numPointsDesc {
    return &numPointsDesc{name: numPointsName, kindType: reflect.Int}
}
type numNodesDesc struct {
    name     string
    value    int
    kindType reflect.Kind
}
func newNumNodesDescRequest(numNodes int) *numNodesDesc {
    return &numNodesDesc{value: numNodes, name: numNodesName, kindType: reflect.Int}
}
func newNumNodesDescResponse() *numNodesDesc {
    return &numNodesDesc{name: numNodesName, kindType: reflect.Int}
}
