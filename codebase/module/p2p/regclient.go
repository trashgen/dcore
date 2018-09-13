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

type regClientModule struct {
    *mediator
}

func newRegClientModule(m *mediator) *regClientModule {
    return &regClientModule{mediator: m}
}

func (this *regClientModule) connect(regListenPort int) {
    this.ThisNodeKey = this.clientModule.SendRequestReg(regListenPort)
    regHosts := this.clientModule.RequestLook(this.ThisNodeKey, this.nodeConfig.MaxPointsCount, this.nodeConfig.MaxP2PConnections)
    for _, regHost := range regHosts {
        log.Printf("RegHost: [%s]", regHost)
        this.createP2PLine(regHost)
    }
}

func (this *regClientModule) createP2PLine(regHost string) {
    conn, err := net.Dial("tcp", regHost)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", regHost)
    }
    if _, err = conn.Write(dcutcp.BuildRequest1013(this.ThisNodeKey)); err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }
    response, err := this.Handle(data, conn)
    if err != nil {
        log.Fatalf("Reg client handler 88 error [%s]: [%s]\n", data, err.Error())
    }
    conn.Write(response)
}

func (this *regClientModule) Handle(data string, conn net.Conn) ([]byte, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil || packetID != dcutcp.RegPacket1013ID {
        return nil, err
    }
    response, err := dcutil.Split1013Response(params)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("regResponseHandler: [%s]", err.Error()))
    }
    if response.Status {
        otherNodeStatus := this.clientModule.RequestCheck(this.ThisNodeKey, response.Target)
        if otherNodeStatus {
            l := newLine(this.mediator, response.Target)
            l.startHost()
            l.startClient(response.Address)
            return dcutcp.BuildRequest88(this.ThisNodeKey, l.address), nil
        } else {
            this.clientModule.SendRequestBan(this.ThisNodeKey, dcutil.RemovePortFromAddressString(conn.RemoteAddr().String()))
            return dcutcp.BuildCommand777(otherNodeStatus), nil
        }
    }
    return nil, errors.New(fmt.Sprintf("regResponseHandler: bad key [%s]", response.Target))
}