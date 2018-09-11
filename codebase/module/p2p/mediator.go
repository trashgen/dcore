package p2p

import (
    "log"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/module/config"
    dchttp "dcore/codebase/module/http/client"
)

// Mediator pattern
type mediator struct {
    ThisNodeKey     string
    lines           map[string]*line
    nodeConfig      *dcconf.NodeConfig
    clientConfig    *dcconf.ClientConfig
    clientModule    *dchttp.HTTPClient
}

func newMediator() *mediator {
    nodeConfig, clientConfig := loadConfigs()
    return &mediator{
        lines           : make(map[string]*line),
        nodeConfig      : nodeConfig,
        clientConfig    : clientConfig,
        clientModule    : dchttp.NewClientModule(clientConfig)}
}

func loadConfigs() (*dcconf.NodeConfig, *dcconf.ClientConfig) {
    c, err := dcutil.LoadJSONConfig(dcconf.NewNodeConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }

    nodeConfig, ok := c.(*dcconf.NodeConfig)
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

    return nodeConfig, clientConfig
}