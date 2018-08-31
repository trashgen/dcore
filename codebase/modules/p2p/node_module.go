package p2p

import (
    "log"
    "net"
    "strings"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

type nodeDesc struct {
    Key         string
    Address     string
    SendConn    net.Conn
    ReceiveConn net.Conn
}

type NodeModule struct {
    *regHostModule
    *regClientModule
    nodes         map[string]*nodeDesc
    config        *dcconf.NodeConfig
    clientConfig  *dcconf.ClientConfig
}

func newNodeDesc(data string) *nodeDesc {
    params := strings.Split(data, "-")
    if len(params) != 2 {
        log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
    }

    return &nodeDesc{Key : params[0], Address : params[1]}
}

func NewNodeModule() *NodeModule {
    c, err := dcutil.LoadJSONConfig(dcconf.NewNodeConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }
    
    config, ok := c.(*dcconf.NodeConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    c, err = dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }

    cmdConfig, ok := c.(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }
    
    c, err = dcutil.LoadJSONConfig(dcconf.NewClientConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }
    
    clientConfig, ok := c.(*dcconf.ClientConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    out := &NodeModule{
        regHostModule   : newRegHostModule(config, cmdConfig, clientConfig),
        regClientModule : newRegClientModule(config),
        nodes           : make(map[string]*nodeDesc),
        config          : config,
        clientConfig    : clientConfig}

    out.parseLookResponse(dchttp.NewClientModule(clientConfig, cmdConfig).RequestLook(1, config.MaxP2PConnections))
    return out
}

func (this *NodeModule) Start() {
    var err error
    this.Key, err = this.startRegHost()
    if err != nil {
        log.Fatal(err.Error())
    }
}