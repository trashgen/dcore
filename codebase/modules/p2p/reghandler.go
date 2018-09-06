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

func (this *regRequestHandler) Run(data string, conn net.Conn) ([]byte, bool, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil {
        return nil, false, err
    }
    switch packetID {
        case dcutcp.RegPacket1013ID():
            return this.handle1013Request(params, conn.RemoteAddr())
        case dcutcp.DeathPacket777ID():
            return this.handle777Request(params)
        case dcutcp.ConfirmPacket88ID():
            return this.handle88Request(params)
    }
    return nil, false, nil
}

func (this *regRequestHandler) handle1013Request(params []string, address net.Addr) ([]byte, bool, error) {
    var err error
    request, err := dcutil.SplitPacket1013RequestParams(params)
    if err != nil {
        return nil, false, err
    }

    var response []byte
    ipOtherNode := dcutil.RemovePortFromAddressString(address.String())
    status := this.clientModule.RequestCheck(request.ThoseNodeKey + "asdqwe")
    if status {
        l := newLine(this.mediator, request.ThoseNodeKey)
        l.startHost()
        response = dcutcp.BuildGoodResponse1013(status, this.ThisNodeKey, l.address)
    } else {
        response = dcutcp.BuildBadResponse1013(status)
        err = errors.New(fmt.Sprintf("Can't reg node with invalid key [%s]\n", request.ThoseNodeKey))
        this.toBlackList <- ipOtherNode
    }

    return response, true, err
}

func (this *regRequestHandler) handle777Request(params []string) ([]byte, bool, error) {
    if command, err := dcutil.SplitCommand777RequestParams(params); ! command.Status || err != nil {
        log.Fatal("To The Death!")
    }
    return nil, false, nil
}

func (this *regRequestHandler) handle88Request(params []string) ([]byte, bool, error) {
    command, err := dcutil.SplitRequest88RequestParams(params)
    if err != nil {
        return nil, false, err
    }
    l := this.lines[command.ThoseNodeKey]
    l.startClient(command.HostAddr)
    return dcutcp.BuildResponse88(this.ThisNodeKey), true, nil
}

///////////////////////////////////////////////////////////////////////////////

type regResponseHandler struct {
    *mediator
}

func newRegResponseHandler(m *mediator) *regResponseHandler {
    return &regResponseHandler{mediator: m}
}

func (this *regResponseHandler) Run(data string, conn net.Conn) ([]byte, bool, error) {
    packetID, params, err := dcutil.SplitPacketIDWithData(data)
    if err != nil {
        return nil, true, err
    }
    switch packetID {
        case dcutcp.RegPacket1013ID():
            return this.handle1013Response(params)
        case dcutcp.ConfirmPacket88ID():
            return this.handle88Response(params)
    }
    return nil, true, nil
}

func (this *regResponseHandler) handle1013Response(params []string) ([]byte, bool, error) {
    response, err := dcutil.SplitPacket1013Response(params)
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

func (this *regResponseHandler) handle88Response(params []string) ([]byte, bool, error) {
    _, err := dcutil.SplitResponse88Params(params)
    if err != nil {
        return nil, false, errors.New(fmt.Sprintf("regResponseHandler: [%s]", err.Error()))
    }
    // decomment for test to start spam message
    //this.lines[response.ThoseNodeKey].Send(response.ThoseNodeKey)
    return nil, false, nil
}