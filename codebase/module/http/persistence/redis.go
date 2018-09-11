package persistence

import "fmt"

type ConnectionID struct {
    IP   string
    Key  string
    Port int
}

func NewConnectionID(key string, ip string, port int) *ConnectionID {
    return &ConnectionID{Key:key, IP:ip, Port:port}
}

func (this ConnectionID) Address() string {
    return fmt.Sprintf("%s:%d", this.IP, this.Port)
}

type RedisModule struct {
    nodes  map[string]*ConnectionID
    points map[string]*ConnectionID // not used now
}

func NewRedisModule() *RedisModule {
    return &RedisModule{
        nodes  : make(map[string]*ConnectionID, 16),
        points : make(map[string]*ConnectionID, 16)}
}

func (this *RedisModule) AddNode(key string, ip string, port int) (out *ConnectionID) {
    out = NewConnectionID(key, ip, port)
    this.nodes[key] = out
    return out
}

func (this *RedisModule) GetNode(key string) (result *ConnectionID, has bool) {
    // LOL I just can't do: "return this.nodes[key]", 'cos "has" not returns by default
    result, has = this.nodes[key]
    return result, has
}

func (this *RedisModule) GetRandomNodes(count int) map[string]*ConnectionID {
    return this.nodes
}

func (this *RedisModule) GetAllNodes() map[string]*ConnectionID {
    return this.nodes
}

func (this *RedisModule) GetAllNodesAsSlice() []*ConnectionID {
    out := make([]*ConnectionID, 0, len(this.nodes))
    for _, v := range this.nodes {
        out = append(out, v)
    }

    return out
}

func (this *RedisModule) RemoveNode(key string) {
    delete(this.nodes, key)
}

func (this *RedisModule) AddPoint(key string, ip string, port int) (out *ConnectionID) {
    out = NewConnectionID(key, ip, port)
    this.points[key] = out
    return out
}

func (this *RedisModule) GetPoint(key string) (result *ConnectionID, has bool) {
    result, has = this.points[key]
    return result, has
}

func (this *RedisModule) GetAllPoints() []*ConnectionID {
    out := make([]*ConnectionID, 0, len(this.points))
    for _, v := range this.points {
        out = append(out, v)
    }
    
    return out
}

func (this *RedisModule) RemovePoint(key string) {
    delete(this.points, key)
}