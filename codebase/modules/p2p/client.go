package p2p

import (
    "log"
    "net"
    "bufio"
    "dcore/codebase/modules/p2p/meta"
)

type client struct {
    *mediator
    hostConn        net.Conn
    regHostAddr     string
    thoseNodeKey    string
    toRequest       chan []byte
    responseHandler meta.ResponseHandler
}

func newClient(m *mediator, thoseNodeKey string, addr string, responseHandler meta.ResponseHandler) *client {
    return &client{
        mediator        : m,
        toRequest       : make(chan []byte),
        regHostAddr     : addr,
        thoseNodeKey    : thoseNodeKey,
        responseHandler : responseHandler}
}

func (this *client) Connect() bool {
    var err error
    this.hostConn, err = net.Dial("tcp", this.regHostAddr)
    if err != nil {
        log.Printf("Can't connect to p2p host [%s]\n", this.regHostAddr)
        return false
    }

    go func() {
        for request := range this.toRequest {
            _, err := this.hostConn.Write(request)
            if err != nil {
                return
            }
            response, err := bufio.NewReader(this.hostConn).ReadString('\n')
            if err != nil {
                return
            }
            err = this.responseHandler.Run(response, this.hostConn)
            if err != nil {
                return
            }
        }
    }()

    return true
}

func (this *client) Send(request string) {
    this.toRequest <- []byte(request)
}