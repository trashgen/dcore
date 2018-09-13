package persistence

import (
    "fmt"
    "log"
    "strconv"
    "strings"
    "github.com/mediocregopher/radix.v3"
)

const (
    setWithKeysName     = "nodekeys"
    numRedisConnections = 4
)

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
    pool    *radix.Pool
    scripts map[string]string
}

func NewRedisModule() *RedisModule {
    p, err := radix.NewPool("tcp", ":6379", numRedisConnections)
    if err != nil {
        log.Fatalln(err)
    }
    if err := p.Do(radix.Cmd(nil, "DEL", setWithKeysName)); err != nil {
        log.Fatalln(err)
    }
    out := &RedisModule{pool: p, scripts: make(map[string]string)}
    out.scripts[getRandomNodesName] = out.compileScript(getRandomNodes)
    return out
}

func (this *RedisModule) AddNode(key string, ip string, port int) *ConnectionID {
    err := this.pool.Do(radix.WithConn("", func (conn radix.Conn) error {
        var err error
        defer func() {
            if err != nil {
                conn.Do(radix.Cmd(nil, "DISCARD"))
                log.Fatalln(err)
            }
        }()
        if err = conn.Do(radix.Cmd(nil, "MULTI")); err != nil {
            return err
        }
        if err = conn.Do(radix.Cmd(nil, "SADD", setWithKeysName, key)); err != nil {
            return err
        }
        if err = conn.Do(radix.Cmd(nil, "HMSET", key, "key", key, "ip", ip, "port", strconv.Itoa(port))); err != nil {
            return err
        }
        if err = conn.Do(radix.Cmd(nil, "EXEC")); err != nil {
            return err
        }
        return nil
    }))
    if err != nil {
        log.Fatalln(err)
    }
    return NewConnectionID(key, ip, port)
}

func (this *RedisModule) NodeExists(key string) bool {
    var i int
    if err := this.pool.Do(radix.Cmd(&i, "HEXISTS", key, "key")); err != nil {
        log.Fatalln(err)
    }
    return i != 0
}

// Методы GetAllNodes/GetRandomNodes - приктически идентичны.
// В GetAllNodes я использую идиому MULTI/EXEC.
func (this *RedisModule) GetAllNodes() []byte {
    sb := strings.Builder{}
    err := this.pool.Do(radix.WithConn("", func (conn radix.Conn) error {
        var err error
        defer func() {
            if err != nil {
                conn.Do(radix.Cmd(nil, "DISCARD"))
                log.Fatalln(err)
            }
        }()
        if err = conn.Do(radix.Cmd(nil, "MULTI")); err != nil {
            return err
        }
        var keys []string
        if err := this.pool.Do(radix.Cmd(&keys, "SMEMBERS", setWithKeysName)); err != nil {
            log.Fatalln(err)
        }
        for _, k := range keys {
            var result map[string]string
            if err := this.pool.Do(radix.Cmd(&result, "HGETALL", k)); err != nil {
                log.Fatalf("HGETALL: [%s]\n", err)
            }
            port, err := strconv.Atoi(result["port"])
            if err != nil {
                log.Fatalln(err)
            }
            sb.WriteString(fmt.Sprintf("%s-%s:%d\t", k, result["ip"], port))
        }
        if err = conn.Do(radix.Cmd(nil, "EXEC")); err != nil {
            return err
        }
        return nil
    }))
    if err != nil {
        log.Fatalf("GetRandomNodes: [%s]\n", err)
    }
    return []byte(sb.String())
}

// Методы GetAllNodes/GetRandomNodes - приктически идентичны.
// В GetRandomNodes я использую Lua stored proc.
func (this *RedisModule) GetRandomNodes(count int) []byte {
    return this.runScriptGetRandomNodes(count)
}

func (this *RedisModule) RemoveNode(key string) {
}

func (this *RedisModule) compileScript(script string) (outSHAKey string) {
    if err := this.pool.Do(radix.Cmd(&outSHAKey, "SCRIPT","LOAD", script)); err != nil {
        log.Fatal(err)
    }
    return outSHAKey
}

func (this *RedisModule) runScriptGetRandomNodes(count int) (result []byte) {
    if err := this.pool.Do(radix.Cmd(&result,"EVALSHA",
                           this.scripts[getRandomNodesName], "0", strconv.Itoa(count))); err != nil {
        log.Fatalln(err)
    }
    return result
}