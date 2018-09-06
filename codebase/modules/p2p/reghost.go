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
                log.Printf("Add Reg Connection [%s]\n", c.RemoteAddr())
                this.createP2PLine(c)
            }(conn)
        }
    }()
}

func (this *regHostModule) onRemoveConnection() {
    go func() {
        for conn := range this.removeConn {
            log.Printf("Remove Reg Connection [%s]\n", conn.RemoteAddr())
        }
    }()
}

func (this *regHostModule) createP2PLine(conn net.Conn) {
    go func() {
        var err error
        for {
            var data string
            data, err = bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                break
            }
            var response []byte
            var hasResponseData bool
            response, hasResponseData, err = this.handler.Run(data, conn)
            if err != nil {
                log.Printf("Reg host handler error [%s]: [%s]\n", data, err.Error())
                break
            }
            if hasResponseData {
                if _, err = conn.Write(response); err != nil {
                    break
                }
            } else {
                break
            }
        }
        this.removeConn <- conn
    }()
}