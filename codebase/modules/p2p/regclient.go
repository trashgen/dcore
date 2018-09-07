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
    if _, _, err = this.Handle(data, conn); err != nil {
        log.Fatalf("Reg client handler 1013 error [%s]: [%s]\n", data, err.Error())
    }
}

func (this *regClientModule) Handle(data string, conn net.Conn) ([]byte, bool, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil || packetID != dcutcp.RegPacket1013ID() {
        return nil, false, err
    }
    response, err := dcutil.Split1013Response(params)
    if err != nil {
        return nil, false, errors.New(fmt.Sprintf("regResponseHandler: [%s]", err.Error()))
    }
    if response.Status {
        otherNodeStatus := this.clientModule.RequestCheck(response.ThoseNodeKey)
        if otherNodeStatus {
            l := newLine(this.mediator, response.ThoseNodeKey)
            l.startHost()
            l.startClient(response.Address)
            return dcutcp.BuildRequest88(this.ThisNodeKey, l.address), true, nil
        } else {
            this.toBlackList <- dcutil.RemovePortFromAddressString(response.Address)
            return dcutcp.BuildCommand777(otherNodeStatus), true, nil
        }
    }
    return nil, false, errors.New(fmt.Sprintf("regResponseHandler: bad key [%s]", response.ThoseNodeKey))
}