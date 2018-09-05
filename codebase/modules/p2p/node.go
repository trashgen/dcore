package p2p

import (
    "log"
    "strings"
    dcutil "dcore/codebase/util"
    "dcore/codebase/modules/persistance"
    "dcore/codebase/modules/p2p/meta"
)

type NodeModule struct {
    *mediator
    *regHostModule
    *regClientModule
}

func NewNodeModule(requestHandler meta.RequestHandler, responseHandler meta.ResponseHandler) *NodeModule {
    mediator := NewMediator(requestHandler, responseHandler)
    out := &NodeModule{mediator: mediator}
    out.regHostModule   = newRegHostModule(mediator)
    out.regClientModule = newRegClientModule(mediator)

    lookResponse := out.clientModule.RequestLook(1, out.nodeConfig.MaxP2PConnections)
    out.handleLookResponse(lookResponse)

    return out
}

func (this *NodeModule) Start() {
    var err error
    this.ThisNodeKey, err = this.startRegHost()
    if err != nil {
        log.Fatal(err.Error())
    }

    go func() {
        postgres := persistance.NewBlackListModule()
        defer postgres.Close()
        for badAddress := range this.toBlackList {
            postgres.Save(badAddress)
        }
    }()
}

func (this *NodeModule) handleLookResponse(data string) {
    data = strings.TrimSuffix(data, "\n")
    values := dcutil.ScanString(data, '\t')
    for _, v := range values {
        params := strings.Split(v, "-")
        if len(params) != 2 {
            log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
        }

        this.clients[params[0]] = newClient(this.mediator, params[0], params[1], this.responseHandler)
    }
}