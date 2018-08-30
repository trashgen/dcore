package p2p

import (
    "fmt"
    "log"
    "net"
    "strings"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

type NodeDesc struct {
    Key         string
    Address     string
    SendConn    net.Conn
    ReceiveConn net.Conn
}

type NodeModule struct {
    *RegHostModule
    *RegClientModule
    Key           string
    nodes         map[string]*NodeDesc
    config        *dcconf.NodeConfig
    clientConfig  *dcconf.ClientConfig
}

func NewNodeDesc(data string) *NodeDesc {
    params := strings.Split(data, ":")
    if len(params) != 3 {
        log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
    }

    return &NodeDesc{
        Key     : params[0],
        Address : fmt.Sprintf("%s:%s", params[1], params[2])}
}

func NewNodeModule(config *dcconf.NodeConfig) *NodeModule {
    clientConfig, ok := dcutil.LoadJSONConfig(dcconf.NewClientConfig(dcconf.NewMetaConfig())).(*dcconf.ClientConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    cmdConfig, ok := dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig())).(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    out := &NodeModule{
        RegHostModule   : NewRegHostModule(config, cmdConfig, clientConfig),
        RegClientModule : NewRegClientModule(config),
        nodes           : make(map[string]*NodeDesc),
        config          : config,
        clientConfig    : clientConfig}

    out.parseLookRequest(dchttp.NewClientModule(clientConfig, cmdConfig).RequestLook(1, config.MaxP2PConnections))
    return out
}

func (this *NodeModule) Start() {
    this.Key = this.startRegHost()
    this.accepting()
}