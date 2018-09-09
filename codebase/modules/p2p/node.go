package p2p

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

func (this *NodeModule) Start() {
    this.connect(this.startRegHost())
}