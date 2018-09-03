package p2p

import (
    "log"
    "net"
    "bufio"
    dcutil "dcore/codebase/util"
    dctcpsrvutil "dcore/codebase/util/tcp/server"
)

type regClientModule struct {
    *mediator
    otherRegHosts []*nodeDesc
}

func newRegClientModule(dataBlock *mediator) *regClientModule {
    return &regClientModule{mediator:dataBlock, otherRegHosts:make([]*nodeDesc, 0, dataBlock.nodeConfig.MaxP2PConnections)}
}

func (this *regClientModule) Connect() {
    for _, nd := range this.otherRegHosts {
        go func(nd *nodeDesc) {
            this.createP2PLine(nd.addr)
        }(nd)
    }
}

func (this *regClientModule) createP2PLine(address string) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", address)
    }

    thisHostAddress := this.StartHost()
    if len(thisHostAddress) == 0 {
        // TODO : обработать, что не могу поднять хост
    }

    request1013 := dctcpsrvutil.BuildPacket1013Request(this.Key)
    log.Printf("request1013 : [%s]\n", request1013)
    _, err = conn.Write(request1013)
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }

    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }

    this.handle1013Response(data, thisHostAddress, conn)
}

func (this *regClientModule) handle1013Response(data string, thisHostAddress string, conn net.Conn) {
    response, err := dcutil.SplitPacket1013Response(data)
    if err != nil {
        log.Printf("Response 1013 bad format [%s]\n", data)
        return
    }

    if response.Status {
        otherNodeStatus := this.clientModule.RequestCheck(response.Key)
        if otherNodeStatus {
            this.StartClient(response.Address)
            _, err := conn.Write(dctcpsrvutil.BuildPacket88(thisHostAddress))
            if err != nil {
                log.Fatalf("Error send Packet 88 to reg host [%s]\n", err.Error())
            }
        } else {
            // TODO : Нужен пакет смерти, который будет убивать другие ноды, если они под подозрением.
            _, err := conn.Write(dctcpsrvutil.BuildPacket777(otherNodeStatus))
            if err != nil {
                log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
            }
            this.toBlackList <- dcutil.RemovePortFromAddressString(response.Address)
        }
    }
}