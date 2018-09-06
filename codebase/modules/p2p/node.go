package p2p

import (
    "dcore/codebase/modules/persistance"
)

type NodeModule struct {
    *mediator
    *regHostModule
    *regClientModule
}

func NewNodeModule() *NodeModule {
    mediator := newMediator()
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