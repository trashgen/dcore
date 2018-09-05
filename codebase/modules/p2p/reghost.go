package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "errors"
    "dcore/codebase/modules/p2p/meta"
)

type regHostModule struct {
    *mediator
    port       int
    handler    meta.RequestHandler
    regHost    net.Listener
    newConn    chan net.Conn
    removeConn chan net.Conn
}

func newRegHostModule(m *mediator) *regHostModule {
    return &regHostModule{
        handler    : newRegRequestHandler(m),
        mediator   : m,
        newConn    : make(chan net.Conn),
        removeConn : make(chan net.Conn)}
}

func (this *regHostModule) startRegHost() (string, error) {
    var err error
    for this.port = this.nodeConfig.MinRegPort; this.port < this.nodeConfig.MaxRegPort; this.port++ {
        if this.regHost, err = net.Listen("tcp", fmt.Sprintf(":%d", this.port)); err == nil {
            break
        }
    }

    if err != nil {
        return "", errors.New(fmt.Sprintf("All reg port are busy [%d..%d]. Do something!\n", this.nodeConfig.MinRegPort, this.nodeConfig.MaxRegPort))
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
            func(c net.Conn) {
                log.Printf("Add Connection [%s]\n", c.RemoteAddr())
                this.createP2PLine(c)
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

func (this *regHostModule) createP2PLine(conn net.Conn) {
    go func() {
        for {
            data, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                break
            }

            response, err := this.handler.Run(data, conn)
            if err != nil {
                break
            }

            if response != nil {
                if _, err = conn.Write(response); err != nil {
                    break
                }
            }
        }
        this.removeConn <- conn
    }()
}