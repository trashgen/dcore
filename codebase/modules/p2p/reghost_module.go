package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

type RegHostModule struct {
    config        *dcconf.NodeConfig
    newConn       chan net.Conn
    regHost       net.Listener
    cmdConfig     *dcconf.HTTPCommands
    removeConn    chan net.Conn
    clientConfig  *dcconf.ClientConfig
}

func NewRegHostModule(config *dcconf.NodeConfig, cmdConfig *dcconf.HTTPCommands, clientConfig  *dcconf.ClientConfig) *RegHostModule {
    return &RegHostModule{
        config       : config,
        newConn      : make(chan net.Conn),
        cmdConfig    : cmdConfig,
        removeConn   : make(chan net.Conn),
        clientConfig : clientConfig}
}

func (this *RegHostModule) startRegHost() string {
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

    this.onNewConnection()
    this.onRemoveConnection()

    return dchttp.NewClientModule(this.clientConfig, this.cmdConfig).RequestReg(fmt.Sprintf("localhost:%d", port))
}

func (this *RegHostModule) accepting() {
    for {
        conn, err := this.regHost.Accept()
        if err != nil {
            log.Fatalf("Can't Accept new connections: [%s]\n", err.Error())
        }
        
        this.newConn <- conn
    }
}

func (this *RegHostModule) onNewConnection() {
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

func (this *RegHostModule) onRemoveConnection() {
    go func() {
        for conn := range this.removeConn {
            log.Printf("Remove Connection [%s]\n", conn.RemoteAddr())
        }
    }()
}

func (this *RegHostModule) processConnRequests(conn net.Conn) {
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

func (this *RegHostModule) handlePacket1013(data string) []byte {
    return []byte(fmt.Sprintf("I am Response from reg client [%s]\n", data))
}