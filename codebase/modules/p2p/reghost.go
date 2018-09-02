package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "errors"
    dcutil "dcore/codebase/util"
    dctcpsrvutil "dcore/codebase/util/tcp/server"
)

type regHostModule struct {
    *mediator
    port       int
    regHost    net.Listener
    newConn    chan net.Conn
    removeConn chan net.Conn
}

func newRegHostModule(dataBlock *mediator) *regHostModule {
    return &regHostModule{
        mediator   : dataBlock,
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
            var err error
            data, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                this.removeConn <- conn
                return
            }

            packetID, params, err := dcutil.SplitPacketIDWithData(data)
            if err != nil {
                this.removeConn <- conn
                return
            }

            switch packetID {
                case dctcpsrvutil.RegPacketID():
                    this.answerTo1013Request(conn, params, conn.RemoteAddr())
                case dctcpsrvutil.DeathPacketID():
                    this.handle777Request(params)
                case dctcpsrvutil.ConfirmPacketID():
                    if this.handle88Request(params) {
                        return
                    }
                    log.Fatal("Error with confirm P2P Connect\n")
            }
        }
    }()
}

func (this *regHostModule) answerTo1013Request(conn net.Conn, params []string, address net.Addr) {
    request, err := dcutil.SplitPacket1013RequestParams(params)
    if err != nil {
        log.Print(err.Error())
        return
    }

    var response []byte
    ipOtherNode := dcutil.RemovePortFromAddressString(address.String())
    status := this.clientModule.RequestCheck(request.Key)
    if status {
        thisHostAddress := this.StartHost()
        if len(thisHostAddress) == 0 {
            // TODO : обработать, что не могу поднять хост
        }
        response = dctcpsrvutil.BuildGoodPacket1013Response(status, this.Key, thisHostAddress)
    } else {
        log.Printf("Can't reg node with invalid key [%s]\n", request.Key)
        response = dctcpsrvutil.BuildBadPacket1013Response(status)
        this.toBlackList <- ipOtherNode
        this.removeConn <- conn
    }

    _, err = conn.Write(response)
    if err != nil {
        this.removeConn <- conn
        return
    }
}

func (this *regHostModule) handle777Request(params []string) {
    if request, err := dcutil.SplitPacket777RequestParams(params); ! request.Status || err != nil {
        log.Fatal("To The Death!")
    }
}

func (this *regHostModule) handle88Request(params []string) bool {
    // TODO : Connect to P2P
    request, err := dcutil.SplitPacket88RequestParams(params)
    if err != nil {
        log.Print("To The Death!")
        return false
    }

    this.StartClient(request.Addr)
    return true
}