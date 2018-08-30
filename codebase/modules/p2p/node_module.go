package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "strings"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcmisc "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
)

type NodeDesc struct {
    Key         string
    Address     string
    SendConn    net.Conn
    ReceiveConn net.Conn
}

type NodeModule struct {
    Key           string
    nodes         map[string]*NodeDesc
    config        *dcconf.NodeConfig
    newConn       chan net.Conn // RegHost struct
    regHost       net.Listener  // RegHost struct
    cmdConfig     *dcconf.HTTPCommands
    removeConn    chan net.Conn // RegHost struct
    clientConfig  *dcconf.ClientConfig
    otherRegHosts []*NodeDesc
}

func NewNodeDesc(data string) *NodeDesc {
    params := strings.Split(data, ":")
    if len(params) != 3 {
        log.Fatalf("Bad Packet1013 Request format [%s]\n", data)
    }

    return &NodeDesc{
        Key     : params[0],
        Address : fmt.Sprintf("%s:%s", params[1], params[2])}
}

func NewNodeModule(config *dcconf.NodeConfig) *NodeModule {
    clientConfig, ok := dcutil.LoadJSONConfig(dcconf.NewClientConfig(dcconf.NewMetaConfig())).(*dcconf.ClientConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    cmdConfig, ok := dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig())).(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    out := &NodeModule{
        nodes         : make(map[string]*NodeDesc),
        config        : config,
        newConn       : make(chan net.Conn),
        cmdConfig     : cmdConfig,
        removeConn    : make(chan net.Conn),
        clientConfig  : clientConfig,
        otherRegHosts : make([]*NodeDesc, 0, config.MaxP2PConnections)}

    out.parseLookRequest(dchttp.NewClientModule(clientConfig, cmdConfig).RequestLook(1, config.MaxP2PConnections))
    return out
}

///////////////////////////////////////////////////////////////////////////////

// RegClient struct
func (this *NodeModule) Connect() {
    for _, nodeDesc := range this.otherRegHosts {
        go func(nd *NodeDesc) {
           this.connectReghost(nd.Address)
        }(nodeDesc)
    }
}

// RegClient struct
func (this *NodeModule) parseLookRequest(data string) {
    data = strings.TrimSuffix(data, "\n")
    scanner := bufio.NewScanner(strings.NewReader(data))
    scanner.Split(dcmisc.SplitterFunc('\t'))
    for scanner.Scan() {
        this.otherRegHosts = append(this.otherRegHosts, NewNodeDesc(scanner.Text()))
    }
}

// RegClient struct
func (this *NodeModule) connectReghost(address string) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", address)
    }

    request := this.createPacket1013()

    _, err = conn.Write(request)
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }
    
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }

    log.Printf("I am Response from reg host [%s]\n", data)
}

// RegClient struct
func (this *NodeModule) createPacket1013() []byte {
    return []byte(fmt.Sprintf("Hi! I am Packet 1013!!!"))
}

///////////////////////////////////////////////////////////////////////////////

// RegHost struct
func (this *NodeModule) StartRegHost() {
    var err error
    var port int
    for port = this.config.MinRegPort; port < this.config.MaxRegPort; port++ {
        log.Printf("Try start to %d\n", port)
        this.regHost, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err != nil {
            log.Printf("Can't start reg host: [%s]\n", err.Error())
        } else {
            break
        }
    }

    if this.regHost == nil {
        log.Fatalf("All reg port are busy [%d..%d]. Do something!\n", this.config.MinRegPort, this.config.MaxRegPort)
    }

    this.Key = dchttp.NewClientModule(this.clientConfig, this.cmdConfig).RequestReg(fmt.Sprintf("localhost:%d", port))
    this.onNewConnection()
    this.onRemoveConnection()

    for {
        conn, err := this.regHost.Accept()
        if err != nil {
            log.Fatalf("Can't Accept new connections: [%s]\n", err.Error())
        }

        this.newConn <- conn
    }
}

// RegHost struct
func (this *NodeModule) onNewConnection() {
    go func() {
        for conn := range this.newConn {
            // IMO best way to work with shared cycle var
            func(c net.Conn) {
                log.Printf("Add Connection [%s]\n", c.RemoteAddr())
                this.processConnRequests(c)
            }(conn)
        }
    }()
}

// RegHost struct
func (this *NodeModule) onRemoveConnection() {
    go func() {
        for conn := range this.removeConn {
            log.Printf("Remove Connection [%s]\n", conn.RemoteAddr())
        }
    }()
}

// RegHost struct
func (this *NodeModule) processConnRequests(conn net.Conn) {
    go func() {
        for {
            data, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                this.removeConn <- conn
                return
                //this.shutdown <- struct{}{}
            }

            response := this.handlePacket1013(data)

            _, err = conn.Write(response)
            if err != nil {
                this.removeConn <- conn
                return
                //this.shutdown <- struct{}{}
            }
        }
    }()
}

// RegHost struct
func (this *NodeModule) handlePacket1013(data string) []byte {
    return []byte(fmt.Sprintf("I am Response from reg client [%s]\n", data))
}