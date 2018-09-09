package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "errors"
    dcutil "dcore/codebase/util"
    dcutcp "dcore/codebase/util/tcp/server"
)

type regHostModule struct {
    *mediator
    newConn     chan net.Conn
    removeConn  chan net.Conn
    regListener net.Listener
}

func newRegHostModule(m *mediator) *regHostModule {
    return &regHostModule{
        mediator          : m,
        newConn           : make(chan net.Conn),
        removeConn        : make(chan net.Conn)}
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
            response, hasResponseData, err = this.Handle(data, conn)
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


func (this *regHostModule) Handle(data string, conn net.Conn) ([]byte, bool, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil {
        return nil, false, err
    }
    switch packetID {
        case dcutcp.RegPacket1013ID:
            return this.handle1013Request(params, conn.RemoteAddr())
        case dcutcp.DeathPacket777ID:
            return this.handle777Command(params)
        case dcutcp.ConfirmPacket88ID:
            return this.handle88Command(params)
    }
    return nil, false, nil
}

func (this *regHostModule) handle1013Request(params []string, address net.Addr) ([]byte, bool, error) {
    var err error
    request, err := dcutil.Split1013RequestParams(params)
    if err != nil {
        return nil, false, err
    }
    
    var response []byte
    ipOtherNode := dcutil.RemovePortFromAddressString(address.String())
    status := this.clientModule.RequestCheck(request.ThoseNodeKey)
    if status {
        l := newLine(this.mediator, request.ThoseNodeKey)
        l.startHost()
        response = dcutcp.BuildGoodResponse1013(status, this.ThisNodeKey, l.address)
    } else {
        response = dcutcp.BuildBadResponse1013(status)
        err = errors.New(fmt.Sprintf("Can't reg node with invalid key [%s]\n", request.ThoseNodeKey))
        // TODO : this.clientModule.RequestBan(request)
        this.toBlackList <- ipOtherNode
    }
    
    return response, true, err
}

func (this *regHostModule) handle777Command(params []string) ([]byte, bool, error) {
    if command, err := dcutil.SplitCommand777Params(params); ! command.Status || err != nil {
        log.Fatal("To The Death!")
    }
    return nil, false, nil
}

func (this *regHostModule) handle88Command(params []string) ([]byte, bool, error) {
    if command, err := dcutil.SplitCommand88Params(params); err == nil {
        l := this.lines[command.ThoseNodeKey]
        l.startClient(command.HostAddr)
    }
    return nil, false, nil
}