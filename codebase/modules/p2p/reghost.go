package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "dcore/codebase/modules/p2p/meta"
)

type regHostModule struct {
    *mediator
    handler     meta.RequestHandler
    newConn     chan net.Conn
    removeConn  chan net.Conn
    regListener net.Listener
}

func newRegHostModule(m *mediator) *regHostModule {
    return &regHostModule{
        handler    : newRegRequestHandler(m),
        mediator   : m,
        newConn    : make(chan net.Conn),
        removeConn : make(chan net.Conn)}
}

func (this *regHostModule) startRegHost() (port int) {
    var err error
    for port = this.nodeConfig.MinRegPort; port < this.nodeConfig.MaxRegPort; port++ {
        if this.regListener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
            break
        }
    }
    if err != nil {
        log.Fatalf("Can't start reg host [%s]\n", err.Error())
    }
    this.onNewConnection()
    this.onRemoveConnection()
    return port
}

func (this *regHostModule) Accepting() {
    for {
        conn, err := this.regListener.Accept()
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
                log.Print(err.Error())
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