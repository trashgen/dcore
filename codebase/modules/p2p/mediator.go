package p2p

import (
    "log"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
    dchttp "dcore/codebase/modules/http/client"
    "dcore/codebase/modules/p2p/meta"
)

// Mediator pattern
type mediator struct {
    ThisNodeKey  string
    hosts        []*host
    clients      map[string]*client
    cmdConfig    *dcconf.HTTPCommands
    nodeConfig   *dcconf.NodeConfig
    clientConfig *dcconf.ClientConfig
    clientModule *dchttp.ClientModule
    toBlackList  chan string
    requestHandler  meta.RequestHandler
    responseHandler meta.ResponseHandler
}

func newMediator(requestHandler meta.RequestHandler, responseHandler meta.ResponseHandler) *mediator {
    nodeConfig, cmdConfig, clientConfig := loadConfigs()
    return &mediator{
        hosts           : make([]*host, 0),
        clients         : make(map[string]*client),
        cmdConfig       : cmdConfig,
        nodeConfig      : nodeConfig,
        toBlackList     : make(chan string, 128),
        clientConfig    : clientConfig,
        clientModule    : dchttp.NewClientModule(clientConfig, cmdConfig),
        requestHandler  : requestHandler,
        responseHandler : responseHandler}
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