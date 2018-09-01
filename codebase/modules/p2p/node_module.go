package p2p

import (
    "log"
    "net"
    "strings"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
    dchttpcli "dcore/codebase/modules/http/client"
    dcpersist "dcore/codebase/modules/persistance"
)

type nodeDesc struct {
    Key         string
    Address     string
    SendConn    net.Conn
    ReceiveConn net.Conn
}

type NodeModule struct {
    *regHostModule
    *regClientModule
    nodes         map[string]*nodeDesc
    config        *dcconf.NodeConfig
    clientConfig  *dcconf.ClientConfig
}

func newNodeDesc(data string) *nodeDesc {
    params := strings.Split(data, "-")
    if len(params) != 2 {
        log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
    }

    return &nodeDesc{Key : params[0], Address : params[1]}
}

func NewNodeModule() *NodeModule {
    c, err := dcutil.LoadJSONConfig(dcconf.NewNodeConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }
    
    config, ok := c.(*dcconf.NodeConfig)
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

    out := &NodeModule{
        regHostModule   : newRegHostModule(config, cmdConfig, clientConfig),
        regClientModule : newRegClientModule(config),
        nodes           : make(map[string]*nodeDesc),
        config          : config,
        clientConfig    : clientConfig}

    out.parseLookResponse(dchttpcli.NewClientModule(clientConfig, cmdConfig).RequestLook(1, config.MaxP2PConnections))
    return out
}

func (this *NodeModule) Start() {
    var err error
    this.Key, err = this.startRegHost()
    if err != nil {
        log.Fatal(err.Error())
    }

    go func() {
        postgres := dcpersist.NewBlackListModule()
        for badAddress := range this.toBlackList {
            postgres.Save(badAddress)
        }
        postgres.Close()
    }()

    go func() {
        for {
            <- this.startP2PHost
            log.Printf("Have to start new P2P Host\n")
            // TODO : Тут работает p2phost_module - создает хост и отдает полный адрес для подключения
            // TODO : Потом начинает работать p2pclient_module - пытается подключится на ip. После коннекта - шлю ответ что готов к обратке.
            // TODO : Учесть момент что коннекта может не быть (NAT) - ограничить по таймауту.
            p2pHostAddress := "P2P Host address"
            this.fromP2PHostModule <- p2pHostAddress
        }
    }()
}