// +build ignore

package p2p

import (
    "fmt"
    "log"
    "bufio"
    "errors"
    "strconv"
    "strings"
    dchttp "dcore/codebase/modules/http"
    dcutil "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
)

type TCPRequestHandler func(data []string) (string, error)

type MessageMapper func(data []string) (int, interface{})

type P2PHandler struct {
    config   *dcconf.TotalConfig
    ssClient *dchttp.SSClient
}

func NewP2PHandler(config *dcconf.TotalConfig, ssClient *dchttp.SSClient) *P2PHandler {
    return &P2PHandler{config:config, ssClient:ssClient}
}

func (this *P2PHandler) Handle(message string) (string, error) {
    log.Printf("Incoming Request: [%s]\n", message)
    out, err := this.parseMessagePacket(message)
    if err != nil {
        return "", err
    }

    log.Printf("Outcoming Response: [%s]\n", out)
    return out, nil
}

func (this *P2PHandler) parseMessagePacket(message string) (string, error) {
    packetID, data, err := this.parseMessageFlat(message)
    if err != nil {
        log.Printf("Error in parseMessagePacket (bad packet): [%s]\nPacket = [%s]\n", err.Error(), message)
        return "", err
    }

    switch packetID {
        case Packet1013ID:
            return this.handlePacket1013(NewPacket1013Request(data)).String(), nil
    }

    log.Fatalf("parseMessagePacket: bad packet ID [%d]\n", packetID)
    return "", nil
}

func (this *P2PHandler) parseMessageFlat(message string) (packetID int, data []string, err error) {
    scanner := bufio.NewScanner(strings.NewReader(message))
    scanner.Split(dcutil.SplitTabs)
    if scanner.Scan() {
        packetID, err = strconv.Atoi(scanner.Text())
        if err != nil {
            return 0, nil, err
        }
    } else {
        return 0, nil, errors.New(fmt.Sprintf("can't read packet ID from message [%s]", message))
    }
    
    for scanner.Scan() {
        data = append(data, scanner.Text())
    }
    
    return packetID, data, nil
}

func (this *P2PHandler) handlePacket1013(packet interface{}) *Packet1013Response {
    packet1013, ok := packet.(*Packet1013Request)
    if ! ok {
        log.Fatalf("executeRealHandler: bad 1013 request")
    }
    
    httpData := this.ssClient.GetRawContent(this.config.SSCommand.Check, packet1013.Key)
    response := NewPacket1013Response(this.ssClient.MapCheck(httpData))
    if response.StateAsBool() {
        // TODO : коннект к хосту регистраций
    }

    return response
}
