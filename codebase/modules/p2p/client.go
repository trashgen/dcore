package p2p

import (
    "fmt"
    "net"
    "errors"
    "log"
)

type client struct {
    *mediator
    outCommand <-chan string
}

func newClient(outCommand <-chan string) *client {
    return &client{outCommand: outCommand}
}

func (this *client) Connect(thoseHostAddr string) error {
    hostConn, err := net.Dial("tcp", thoseHostAddr)
    if err != nil {
        return errors.New(fmt.Sprintf("Can't connect to p2p host [%s]: [%s]\n", thoseHostAddr, err.Error()))
    }

    this.handleConn(hostConn)
    return nil
}

func (this *client) handleConn(conn net.Conn) {
    go func() {
        for outCommand := range this.outCommand {
            request, hasRequest := doSomeOut(outCommand)
            if hasRequest {
                // decomment for test that has echo message
                // log.Println(request)
                if _, err := conn.Write([]byte(request)); err != nil {
                    log.Println("Client: Error on send data to host")
                }
            }
        }
    }()
}

func doSomeOut(outCommand string) (string, bool) {
    return fmt.Sprintf("Request for: [%s]\n", outCommand), true
}