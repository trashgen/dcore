package p2p

import "net"

type processor struct {
    Request  string
    Response []byte
    sendConn chan net.Conn
    recvConn chan net.Conn
}

func NewProcessor() *processor {
    return &processor{sendConn:make(chan net.Conn), recvConn:make(chan net.Conn)}
}

func (this *processor) Handle(request string) (response []byte) {
    response = make([]byte, 0)
    return response
}