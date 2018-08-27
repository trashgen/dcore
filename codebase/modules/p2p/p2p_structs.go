package p2p

import (
    "fmt"
    "log"
    "strconv"
    dchttp "dcore/codebase/modules/http"
)

const (
    Packet1013ID        = 1013
    packet1013NumParams = 4
)

type Packet1013Request struct {
    Key       string
    regIP     string
    regPort   int
    SSAddress string
}

type Packet1013Response struct {
    SSCheckState string
}

func NewPacket1013Request(data []string) *Packet1013Request {
    if len(data) != packet1013NumParams {
        // TODO : реализовать логику валидации пакетов (по количеству параметров)
        log.Fatalf("mapToPacket1013: bad params count [%d]\n", len(data))
    }

    port, err := strconv.Atoi(data[3])
    if err != nil {
        log.Fatalf("mapToPacket1013: port not int [%#v]\n", data[3])
    }

    return &Packet1013Request{
        Key       : data[0],
        regIP     : data[2],
        regPort   : port,
        SSAddress : data[1]}
}

func NewPacket1013Response(httpResponse *dchttp.ResponseCheck) *Packet1013Response {
    return &Packet1013Response{SSCheckState:strconv.FormatBool(httpResponse.OpResult)}
}

func (this *Packet1013Response) String() string {
    return fmt.Sprintf("%d\t%s\n", Packet1013ID, this.SSCheckState)
}

func (this *Packet1013Response) StateAsBool() bool {
    out, err := strconv.ParseBool(this.SSCheckState)
    if err != nil {
        log.Fatal(err.Error())
    }

    return out
}