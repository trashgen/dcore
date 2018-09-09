package config

import (
    "fmt"
    "net/url"
    "reflect"
    "strconv"
)

const (
    HTTPParamKey       = "key"
    HTTPParamPort      = "port"
    HTTPParamCount     = "count"
    HTTPParamNodes     = "nodes"
    HTTPParamPoints    = "points"
    HTTPParamTarget    = "target"
    metaConfigFileName = "meta.cfg"
)

type MetaConfig struct {
    NodeConfigFileName   string
    PointConfigFileName  string
    ClientConfigFileName string
    HTTPCommandsFileName string
}

func NewMetaConfig() *MetaConfig {
    return &MetaConfig{
        NodeConfigFileName   : "nodeconfig.cfg",
        PointConfigFileName  : "pointconfig.cfg",
        ClientConfigFileName : "clientconfig.cfg",
        HTTPCommandsFileName : "httpcmdconfig.cfg"}
}

func (this *MetaConfig) String() string {
    return metaConfigFileName
}

///////////////////////////////////////////////////////////////////////////////

type NodeConfig struct {
    MinRegPort        int
    MaxRegPort        int
    MinP2PPort        int
    MaxP2PPort        int
    MaxPointsCount    int
    MaxP2PConnections int
    fileName          string
}

func NewNodeConfig(meta *MetaConfig) *NodeConfig {
    return &NodeConfig{
        MinRegPort        : 33333,
        MaxRegPort        : 33366,
        MinP2PPort        : 51111,
        MaxP2PPort        : 51222,
        MaxPointsCount    : 1,
        MaxP2PConnections : 16,
        fileName          : meta.NodeConfigFileName}
}

func (this *NodeConfig) String() string {
    return this.fileName
}

///////////////////////////////////////////////////////////////////////////////

type CommandDesc struct {
    Name   string
    Params map[string]reflect.Kind
}

type HTTPCommands struct {
    Ban      *CommandDesc
    Reg      *CommandDesc
    Look     *CommandDesc
    Root     *CommandDesc
    Check    *CommandDesc
    Points   *CommandDesc
    Remove   *CommandDesc
    fileName string
    allCommands []*CommandDesc
}

func NewHTTPCommands(meta *MetaConfig) *HTTPCommands {
    out := &HTTPCommands{
        Root         : &CommandDesc{Name:"",       Params:map[string]reflect.Kind{}},
        Reg          : &CommandDesc{Name:"reg",    Params:map[string]reflect.Kind{"port"  : reflect.Int}},
        Points       : &CommandDesc{Name:"points", Params:map[string]reflect.Kind{"count" : reflect.Int}},
        Look         : &CommandDesc{Name:"look",   Params:map[string]reflect.Kind{"points": reflect.Int,    "nodes" : reflect.Int}},
        Ban          : &CommandDesc{Name:"ban",    Params:map[string]reflect.Kind{"key"   : reflect.String, "target": reflect.String}},
        Check        : &CommandDesc{Name:"check",  Params:map[string]reflect.Kind{"key"   : reflect.String, "target": reflect.String}},
        Remove       : &CommandDesc{Name:"remove", Params:map[string]reflect.Kind{"key"   : reflect.String, "target": reflect.String}},
        fileName     : meta.HTTPCommandsFileName}
    out.linkAllCommands()
    return out
}

func (this HTTPCommands) IsValidRequest(cmd string, params url.Values) bool {
    for _, c := range this.allCommands {
        if c.Name == cmd && len(params) == len(c.Params) {
            for k := range params {
                paramType, ok := c.Params[k]
                if ! ok {
                    return false
                }
                if paramType == reflect.Int {
                    if _, err := strconv.Atoi(params.Get(k)); err != nil {
                        return false
                    }
                }
            }
            return true
        }
    }
    return false
}

func (this HTTPCommands) String() string {
    return this.fileName
}

func (this HTTPCommands) linkAllCommands() {
    // As fast as can
    this.allCommands = make([]*CommandDesc, 7)
    this.allCommands[0] = this.Ban
    this.allCommands[1] = this.Reg
    this.allCommands[2] = this.Look
    this.allCommands[3] = this.Root
    this.allCommands[4] = this.Check
    this.allCommands[5] = this.Points
    this.allCommands[6] = this.Remove
}

///////////////////////////////////////////////////////////////////////////////

type PointConfig struct {
    fileName     string
    ListenPort   int
    SecretPhrase string
}

func NewPointConfig(meta *MetaConfig) *PointConfig {
    return &PointConfig{
        ListenPort   : 30001,
        SecretPhrase : "operation cwal (C) Starcraft",
        fileName     : meta.PointConfigFileName}
}

func (this *PointConfig) FormattedListenPort() string {
    return fmt.Sprintf(":%d", this.ListenPort)
}

func (this *PointConfig) String() string {
    return this.fileName
}

///////////////////////////////////////////////////////////////////////////////

type ClientConfig struct {
    EntryPoints []string
    fileName    string
}

func NewClientConfig(meta *MetaConfig) *ClientConfig {
    points := make([]string, 0)
    points = append(points, "http://localhost:30001")

    return &ClientConfig{EntryPoints : points, fileName : meta.ClientConfigFileName}
}

func (this *ClientConfig) String() string {
    return this.fileName
}