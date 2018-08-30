package config

import "fmt"

const metaConfigFileName = "meta.cfg"

type MetaConfig struct {
    NodeConfigFileName   string
    PointConfigFileName  string
    ClientConfigFileName string
}

func NewMetaConfig() *MetaConfig {
    return &MetaConfig{
        NodeConfigFileName   : "nodeconfig.cfg",
        PointConfigFileName  : "pointconfig.cfg",
        ClientConfigFileName : "clientconfig.cfg"}
}

func (this *MetaConfig) String() string {
    return metaConfigFileName
}

///////////////////////////////////////////////////////////////////////////////

type NodeConfig struct {
    MinRegPort        int
    MaxRegPort        int
    MaxP2PConnections int
    fileName          string
}

func NewNodeConfig(meta *MetaConfig) *NodeConfig {
    return &NodeConfig{
        MinRegPort        : 33333,
        MaxRegPort        : 33366,
        MaxP2PConnections : 16,
        fileName          : meta.NodeConfigFileName}
}

func (this *NodeConfig) String() string {
    return this.fileName
}

///////////////////////////////////////////////////////////////////////////////

// TODO : Разместить в отдельном файле чтоб все если с одной миски!!!!
type CommandDesc struct {
    Name  string
    Param string
}

type PointConfig struct {
    Reg          CommandDesc
    Look         CommandDesc
    Root         CommandDesc
    Check        CommandDesc
    Points       CommandDesc
    Remove       CommandDesc
    fileName     string
    ListenPort   int
    SecretPhrase string
}

func NewPointConfig(meta *MetaConfig) *PointConfig {
    return &PointConfig{
        Reg          : CommandDesc{Name:"reg", Param:"address"},
        Look         : CommandDesc{Name:"look", Param:"count"},
        Root         : CommandDesc{Name:""},
        Check        : CommandDesc{Name:"check", Param:"key"},
        Points       : CommandDesc{Name:"points", Param:"count"},
        Remove       : CommandDesc{Name:"remove", Param:"key"},
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
    Reg         CommandDesc
    Look        CommandDesc
    Root        CommandDesc
    Check       CommandDesc
    Points      CommandDesc
    Remove      CommandDesc
    EntryPoints []string
    fileName    string
}

func NewClientConfig(meta *MetaConfig) *ClientConfig {
    points := make([]string, 0)
    points = append(points, "http://localhost:30001")
    
    return &ClientConfig{
        Reg         : CommandDesc{Name:"reg", Param:"address"},
        Look        : CommandDesc{Name:"look", Param:"count"},
        Root        : CommandDesc{Name:""},
        Check       : CommandDesc{Name:"check", Param:"key"},
        Points      : CommandDesc{Name:"points", Param:"count"},
        Remove      : CommandDesc{Name:"remove", Param:"key"},
        EntryPoints : points,
        fileName    : meta.ClientConfigFileName}
}

func (this *ClientConfig) String() string {
    return this.fileName
}