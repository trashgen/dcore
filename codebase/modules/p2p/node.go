package p2p

import (
    "log"
    "strings"
    dcutil "dcore/codebase/util"
    "dcore/codebase/modules/persistance"
)

type nodeDesc struct {
    key  string
    addr string
}

type NodeModule struct {
    *mediator
    *regHostModule
    *regClientModule
}

func newNodeDesc(data string) *nodeDesc {
    params := strings.Split(data, "-")
    if len(params) != 2 {
        log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
    }

    return &nodeDesc{key: params[0], addr: params[1]}
}

func NewNodeModule() *NodeModule {
    mediator := NewMediator()
    out := &NodeModule{mediator:mediator}
    out.regHostModule   = newRegHostModule(mediator)
    out.regClientModule = newRegClientModule(mediator)

    lookResponse := out.clientModule.RequestLook(1, out.nodeConfig.MaxP2PConnections)
    out.parseLookResponse(lookResponse)

    return out
}

func (this *NodeModule) Start() {
    var err error
    this.Key, err = this.startRegHost()
    if err != nil {
        log.Fatal(err.Error())
    }

    go func() {
        postgres := persistance.NewBlackListModule()
        for badAddress := range this.toBlackList {
            postgres.Save(badAddress)
        }
        postgres.Close()
    }()

    go func() {
        for {
            <- this.startP2PHost
            host := NewP2PHost()
            address, status := host.Start(this.nodeConfig)
            if ! status {
                log.Printf("P2P Host : no more empty ports. Check nodeconfig.cfg\n")
            }
            this.fromP2PHostModule <- address
        }
    }()
    
    go func() {
        for {
            hostAddress := <- this.startP2PClient
            client := NewP2PClient()
            client.Connect(hostAddress)
            this.fromP2PClientModule <- struct{}{}
        }
    }()
}

func (this *regClientModule) parseLookResponse(data string) {
    data = strings.TrimSuffix(data, "\n")
    values := dcutil.ScanString(data, '\t')
    for _, nd := range values {
        this.otherRegHosts = append(this.otherRegHosts, newNodeDesc(nd))
    }
}