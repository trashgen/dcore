package p2p

import (
    "log"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
    dchttp "dcore/codebase/modules/http/client"
)

// Mediator pattern
type mediator struct {
    ThisNodeKey     string
    lines           map[string]*line
    toBlackList     chan string
    cmdConfig       *dcconf.HTTPCommands
    nodeConfig      *dcconf.NodeConfig
    clientConfig    *dcconf.ClientConfig
    clientModule    *dchttp.HTTPClient
}

func newMediator() *mediator {
    nodeConfig, cmdConfig, clientConfig := loadConfigs()
    return &mediator{
        lines           : make(map[string]*line),
        cmdConfig       : cmdConfig,
        nodeConfig      : nodeConfig,
        toBlackList     : make(chan string, 128),
        clientConfig    : clientConfig,
        clientModule    : dchttp.NewClientModule(clientConfig, cmdConfig)}
}

func loadConfigs() (*dcconf.NodeConfig, *dcconf.HTTPCommands, *dcconf.ClientConfig) {
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

    return nodeConfig, cmdConfig, clientConfig
}