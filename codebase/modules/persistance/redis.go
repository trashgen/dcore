package persistance

import (
    dchttputil "dcore/codebase/util/http"
)

type RedisModule struct {
    nodes  map[string]*dchttputil.ConnectionID
    points map[string]*dchttputil.ConnectionID // not used now
}

func NewRedisModule() *RedisModule {
    return &RedisModule{
        nodes  : make(map[string]*dchttputil.ConnectionID, 16),
        points : make(map[string]*dchttputil.ConnectionID, 16)}
}

func (this *RedisModule) AddNode(key string, ip string, port int) (out *dchttputil.ConnectionID) {
    out = dchttputil.NewConnectionID(key, ip, port)
    this.nodes[key] = out
    return out
}

func (this *RedisModule) GetNode(key string) (result *dchttputil.ConnectionID, has bool) {
    result, has = this.nodes[key]
    return result, has
}

func (this *RedisModule) GetAllNodes() map[string]*dchttputil.ConnectionID {
    return this.nodes
}

func (this *RedisModule) GetAllNodesAsSlice() []*dchttputil.ConnectionID {
    out := make([]*dchttputil.ConnectionID, 0, len(this.nodes))
    for _, v := range this.nodes {
        out = append(out, v)
    }

    return out
}

func (this *RedisModule) RemoveNode(key string) {
    delete(this.nodes, key)
}

func (this *RedisModule) AddPoint(key string, ip string, port int) (out *dchttputil.ConnectionID) {
    out = dchttputil.NewConnectionID(key, ip, port)
    this.points[key] = out
    return out
}

func (this *RedisModule) GetPoint(key string) (result *dchttputil.ConnectionID, has bool) {
    result, has = this.points[key]
    return result, has
}

func (this *RedisModule) GetAllPoints() []*dchttputil.ConnectionID {
    out := make([]*dchttputil.ConnectionID, 0, len(this.points))
    for _, v := range this.points {
        out = append(out, v)
    }
    
    return out
}

func (this *RedisModule) RemovePoint(key string) {
    delete(this.points, key)
}