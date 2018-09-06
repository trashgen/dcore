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
    var err error
    regHosts := this.clientModule.RequestLook(1, this.nodeConfig.MaxP2PConnections)
    this.ThisNodeKey, err = this.clientModule.RequestReg(regListenPort)
    log.Printf("My reg key is [%s]\n", this.ThisNodeKey)
    if err != nil {
        log.Fatalf("Can't register at Point [%s]\n", err.Error())
    }
    for _, regHost := range regHosts {
        regHost := regHost
        this.createP2PLine(regHost)
    }
}

// TODO : зарефакторить этот гвонокод!!!!!!!!!
func (this *regClientModule) createP2PLine(regHost string) {
    conn, err := net.Dial("tcp", regHost)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", regHost)
    }

    request1013 := dcutcp.BuildRequest1013(this.ThisNodeKey)
    _, err = conn.Write(request1013)
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }

    request88, _, err := this.handler.Run(data, conn)
    if err != nil {
        log.Printf("Reg client handler 1013 error [%s]: [%s]\n", data, err.Error())
    }
    
    _, err = conn.Write(request88)
    if err != nil {
        log.Fatalf("Error send Packet 88 to reg host [%s]\n", err.Error())
    }
    
    data, err = bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 88 from reg host [%s]\n", err.Error())
    }

    if _, _, err = this.handler.Run(data, conn); err != nil {
        log.Printf("Reg client handler 88 error [%s]: [%s]\n", data, err.Error())
    }
}