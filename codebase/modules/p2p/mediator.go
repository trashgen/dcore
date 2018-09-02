package p2p

import (
    "log"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
    dchttpcli "dcore/codebase/modules/http/client"
)

// Mediator pattern
type mediator struct {
    Key   string
    nodes map[string]*nodeDesc

    cmdConfig    *dcconf.HTTPCommands
    nodeConfig   *dcconf.NodeConfig
    clientConfig *dcconf.ClientConfig

    clientModule *dchttpcli.ClientModule

    toBlackList         chan string
    startP2PHost        chan struct{}
    startP2PClient      chan string
    fromP2PHostModule   chan string
    fromP2PClientModule chan struct{}
}

func NewMediator() *mediator {
    c, err := dcutil.LoadJSONConfig(dcconf.NewNodeConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }

    nodeConfig, ok := c.(*dcconf.NodeConfig)
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

    return &mediator{
        nodes               : make(map[string]*nodeDesc),
        cmdConfig           : cmdConfig,
        nodeConfig          : nodeConfig,
        clientConfig        : clientConfig,
        clientModule        : dchttpcli.NewClientModule(clientConfig, cmdConfig),
        toBlackList         : make(chan string, 128),
        startP2PHost        : make(chan struct{}),
        startP2PClient      : make(chan string),
        fromP2PHostModule   : make(chan string),
        fromP2PClientModule : make(chan struct{})}
}

func (this *mediator) StartHost() string {
    this.startP2PHost <- struct{}{}
    // TODO : add to nodes
    return <- this.fromP2PHostModule
}

func (this *mediator) StartClient(onHost string) {
    this.startP2PClient <- onHost
    // TODO : update nodes record
    <- this.fromP2PClientModule
}