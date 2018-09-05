package p2p

import (
    "log"
    "net"
    "bufio"
    "dcore/codebase/modules/p2p/meta"
    dcutcp "dcore/codebase/util/tcp/server"
)

type regClientModule struct {
    *mediator
    handler meta.ResponseHandler
}

func newRegClientModule(m *mediator) *regClientModule {
    return &regClientModule{mediator: m, handler: newRegResponseHandler(m)}
}

func (this *regClientModule) Connect(regListenPort int) {
    regHosts := this.clientModule.RequestLook(1, this.nodeConfig.MaxP2PConnections)
    var err error
    this.ThisNodeKey, err = this.clientModule.RequestReg(regListenPort)
    if err != nil {
        log.Fatalf("Can't register at Point [%s]\n", err.Error())
    }
    for _, regHost := range regHosts {
        regHost := regHost
        this.createP2PLine(regHost)
    }
}

func (this *regClientModule) createP2PLine(regHost string) {
    conn, err := net.Dial("tcp", regHost)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", regHost)
    }

    request1013 := dcutcp.BuildPacket1013Request(this.ThisNodeKey)
    _, err = conn.Write(request1013)
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }

    this.handler.Run(data, conn)
}