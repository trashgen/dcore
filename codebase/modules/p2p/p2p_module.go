// +build ignore

package p2p

import (
    "log"
    "net"
    dctcp "dcore/codebase/modules/tcp"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

type NodeModule struct {
    ID             string
    config         *dcconf.TotalConfig
    regModule      *RegModule
    p2pModules     []*dctcp.TCPModule
    regTCPModule   *dctcp.TCPModule
    connectedNodes map[string]*dchttp.NodeID
}

func NewNodeModule(config *dcconf.TotalConfig) *NodeModule {
    ssClient := dchttp.NewSSClient(config)
    regTCPModule := dctcp.NewTCPModule(config, NewP2PHandler(config, ssClient), 16)
    regTCPModule.OnNewConnection = func(conn net.Conn) {
        log.Printf("New connection on reg host [%s]\n", conn.RemoteAddr().String())
    }
    regTCPModule.OnRemoveConnection = func(conn net.Conn) {
        log.Printf("Remove connection on reg host [%s]\n", conn.RemoteAddr().String())
    }
    regTCPModule.OnShutdown = func() {
        log.Print("Shutdown...\n")
    }

    return &NodeModule {
        config         : config,
        regModule      : NewRegModule(config, ssClient),
        p2pModules     : make([]*dctcp.TCPModule, 0, config.Node.RequestActiveNodesCount),
        regTCPModule   : regTCPModule,
        connectedNodes : make(map[string]*dchttp.NodeID, config.Node.RequestActiveNodesCount)}
}

func (this *NodeModule) StartBuildingP2P() {
    go func() {
        this.regTCPModule.StartHost(true, 0)
    }()

    go func() {
        this.regModule.Execute()
    }()
}