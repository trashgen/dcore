package config

import "fmt"

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
}

func NewMetaConfig() *MetaConfig {
	return &MetaConfig{
		NodeConfigFileName:   "nodeconfig.cfg",
		PointConfigFileName:  "pointconfig.cfg",
		ClientConfigFileName: "clientconfig.cfg"}
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
		MinRegPort:        33333,
		MaxRegPort:        33366,
		MinP2PPort:        51111,
		MaxP2PPort:        51222,
		MaxPointsCount:    1,
		MaxP2PConnections: 16,
		fileName:          meta.NodeConfigFileName}
}

func (this *NodeConfig) String() string {
	return this.fileName
}

///////////////////////////////////////////////////////////////////////////////

type PointConfig struct {
	fileName     string
	ListenPort   int
	SecretPhrase string
}

func NewPointConfig(meta *MetaConfig) *PointConfig {
	return &PointConfig{
		ListenPort:   30001,
		SecretPhrase: "operation cwal (C) Starcraft",
		fileName:     meta.PointConfigFileName}
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

	return &ClientConfig{EntryPoints: points, fileName: meta.ClientConfigFileName}
}

func (this *ClientConfig) String() string {
	return this.fileName
}
