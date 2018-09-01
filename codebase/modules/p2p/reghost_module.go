package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "errors"
    "strings"
    dcutil "dcore/codebase/util"
    dcmisc "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
    dctcpsrvutil "dcore/codebase/util/tcp/server"
    dchttpcli "dcore/codebase/modules/http/client"
)

type regHostModule struct {
    port         int
    config       *dcconf.NodeConfig
    newConn      chan net.Conn
    regHost      net.Listener
    cmdConfig    *dcconf.HTTPCommands
    removeConn   chan net.Conn
    clientConfig *dcconf.ClientConfig
    clientModule *dchttpcli.ClientModule

    toBlackList       chan string
    startP2PHost      chan struct{}
    fromP2PHostModule chan string
}

func newRegHostModule(config *dcconf.NodeConfig, cmdConfig *dcconf.HTTPCommands, clientConfig  *dcconf.ClientConfig) *regHostModule {
    return &regHostModule{
        config       : config,
        newConn      : make(chan net.Conn),
        cmdConfig    : cmdConfig,
        removeConn   : make(chan net.Conn),
        clientConfig : clientConfig,
        clientModule : dchttpcli.NewClientModule(clientConfig, cmdConfig),

        toBlackList       : make(chan string, 128),
        startP2PHost      : make(chan struct{}),
        fromP2PHostModule : make(chan string)}
}

func (this *regHostModule) startRegHost() (string, error) {
    var err error
    for this.port = this.config.MinRegPort; this.port < this.config.MaxRegPort; this.port++ {
        if this.regHost, err = net.Listen("tcp", fmt.Sprintf(":%d", this.port)); err == nil {
            break
        }
    }

    if err != nil {
        return "", errors.New(fmt.Sprintf("All reg port are busy [%d..%d]. Do something!\n", this.config.MinRegPort, this.config.MaxRegPort))
    }

    this.onNewConnection()
    this.onRemoveConnection()

    return this.clientModule.RequestReg(this.port)
}

func (this *regHostModule) Accepting() {
    for {
        conn, err := this.regHost.Accept()
        if err != nil {
            log.Fatalf("Can't Accept new connections: [%s]\n", err.Error())
        }
        
        this.newConn <- conn
    }
}

func (this *regHostModule) onNewConnection() {
    go func() {
        for conn := range this.newConn {
            go func(c net.Conn) {
                log.Printf("Add Connection [%s]\n", c.RemoteAddr())
                this.processPacket1013(c)
            }(conn)
        }
    }()
}

func (this *regHostModule) onRemoveConnection() {
    go func() {
        for conn := range this.removeConn {
            log.Printf("Remove Connection [%s]\n", conn.RemoteAddr())
        }
    }()
}

func (this *regHostModule) processPacket1013(conn net.Conn) {
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        this.removeConn <- conn
        return
    }

    response, key, err := this.handlePacket1013(data, conn.RemoteAddr())
    if err != nil {
        log.Printf("Can't reg node with invalid key [%s]\n", key)
        this.removeConn <- conn

        return
    }

    _, err = conn.Write(response)
    if err != nil {
        this.removeConn <- conn
        return
    }
}

func (this *regHostModule) handlePacket1013(data string, address net.Addr) ([]byte, string, error) {
    request, err := dcmisc.SplitPacket1013(data)
    if err != nil {
        return nil, "", err
    }

    response := this.clientModule.RequestCheck(request.Key)
    if strings.TrimSuffix(response, "\n") == "true" {
        this.startP2PHost <- struct{}{}
        hostAddress := <- this.fromP2PHostModule
        log.Printf("New P2P Host address is [%s]\n", hostAddress)
        response = dctcpsrvutil.BuildGoodPacket1013Response(hostAddress)
    } else {
        log.Printf("Address [%s] to black list\n", address)
        this.toBlackList <- dcutil.RemovePortFromAddressString(address.String())
    }

    return []byte(response), request.Key, nil
}