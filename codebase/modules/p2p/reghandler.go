package p2p

import (
    "fmt"
    "log"
    "net"
    "errors"
    dcutil "dcore/codebase/util"
    dcutcp "dcore/codebase/util/tcp/server"
)

type regRequestHandler struct {
    *mediator
}

func newRegRequestHandler(m *mediator) *regRequestHandler {
    return &regRequestHandler{mediator: m}
}

func (this *regRequestHandler) Run(data string, conn net.Conn) ([]byte, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil {
        return nil, err
    }
    switch packetID {
        case dcutcp.RegPacketID():
            return this.handle1013Request(conn, params, conn.RemoteAddr())
        case dcutcp.DeathPacketID():
            return this.handle777Request(params)
        case dcutcp.ConfirmPacketID():
            if err := this.handle88Request(params); err != nil {
                log.Fatalf("Error with confirm P2P Connect [%s]\n", err.Error())
            }
    }
    return nil, nil
}

func (this *regRequestHandler) handle1013Request(conn net.Conn, params []string, address net.Addr) ([]byte, error) {
    var err error
    request, err := dcutil.SplitPacket1013RequestParams(params)
    if err != nil {
        return nil, err
    }

    var response []byte
    ipOtherNode := dcutil.RemovePortFromAddressString(address.String())
    status := this.clientModule.RequestCheck(request.ThoseNodeKey)
    if status {
        thisHostAddress := startHost(this.ThisNodeKey, this.mediator)
        if len(thisHostAddress) == 0 {
            // TODO : обработать, что не могу поднять хост
        }
        response = dcutcp.BuildGoodPacket1013Response(status, this.ThisNodeKey, thisHostAddress)
    } else {
        response = dcutcp.BuildBadPacket1013Response(status)
        err = errors.New(fmt.Sprintf("Can't reg node with invalid key [%s]\n", request.ThoseNodeKey))
        this.toBlackList <- ipOtherNode
    }

    return response, err
}

func (this *regRequestHandler) handle777Request(params []string) ([]byte, error) {
    if request, err := dcutil.SplitPacket777RequestParams(params); ! request.Status || err != nil {
        log.Fatal("To The Death!")
    }
    return nil, nil
}

func (this *regRequestHandler) handle88Request(params []string) error {
    request, err := dcutil.SplitPacket88RequestParams(params)
    if err != nil {
        return err
    }
    client := startClient(request.ThoseNodeKey, request.HostAddr, this.mediator)
    client.Send("First fuckin' message !!!!\n")
    return nil
}

///////////////////////////////////////////////////////////////////////////////

type regResponseHandler struct {
    *mediator
}

func newRegResponseHandler(m *mediator) *regResponseHandler {
    return &regResponseHandler{mediator: m}
}

func (this *regResponseHandler) Run(data string, conn net.Conn) error {
    response, err := dcutil.SplitPacket1013Response(data)
    if err != nil {
        return errors.New(fmt.Sprintf("regResponseHandler: [%s]", err.Error()))
    }

    if response.Status {
        otherNodeStatus := this.clientModule.RequestCheck(response.ThoseNodeKey)
        if otherNodeStatus {
            thisHostAddress := startHost(this.ThisNodeKey, this.mediator)
            startClient(response.ThoseNodeKey, response.Address, this.mediator)
            _, err := conn.Write(dcutcp.BuildPacket88(this.ThisNodeKey, thisHostAddress))
            if err != nil {
                log.Fatalf("Error send Packet 88 to reg host [%s]\n", err.Error())
            }
        } else {
            _, err := conn.Write(dcutcp.BuildPacket777(otherNodeStatus))
            if err != nil {
                log.Fatalf("Error send Packet 777 to reg host [%s]\n", err.Error())
            }
            this.toBlackList <- dcutil.RemovePortFromAddressString(response.Address)
        }
        return nil
    }

    return errors.New(fmt.Sprintf("regResponseHandler: bad key [%s]", response.ThoseNodeKey))
}

func startHost(thisNodeKey string, m *mediator) (address string) {
    var err error
    host := newHost(thisNodeKey, m.requestHandler)
    address, err = host.Start(m.nodeConfig.MinP2PPort, m.nodeConfig.MaxP2PPort)
    if err != nil {
        log.Print(err.Error())
    }
    m.hosts = append(m.hosts, host)
    return address
}

func startClient(thoseNodeKey string, hostAddr string, m *mediator) *client {
    client := newClient(m, thoseNodeKey, hostAddr, m.responseHandler)
    m.clients[thoseNodeKey] = client
    if ! client.Connect() {
        // TODO : обработать невозможность прямого подключения
    }
    return client
}