package p2p

import (
    "dcore/codebase/modules/p2p/meta"
    "dcore/codebase/modules/persistance"
)

type NodeModule struct {
    *mediator
    *regHostModule
    *regClientModule
}

func NewNodeModule(requestHandler meta.RequestHandler, responseHandler meta.ResponseHandler) *NodeModule {
    mediator := newMediator(requestHandler, responseHandler)
    return &NodeModule{
        mediator        : mediator,
        regHostModule   : newRegHostModule(mediator),
        regClientModule : newRegClientModule(mediator)}
}

func (this *NodeModule) Start() (regListenPort int) {
    regListenPort = this.startRegHost()
    go func() {
        postgres := persistance.NewBlackListModule()
        defer postgres.Close()
        for badAddress := range this.toBlackList {
            badAddress := badAddress
            postgres.Save(badAddress)
        }
    }()
    return regListenPort
}