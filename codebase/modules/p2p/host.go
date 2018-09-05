package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "dcore/codebase/modules/p2p/meta"
)

type host struct {
    key            string
    requestHandler meta.RequestHandler
}

func newHost(key string, requestHandler meta.RequestHandler) *host {
    return &host{key: key, requestHandler: requestHandler}
}

func (this *host) Start(minPort int, maxPort int) (string, error) {
    var port int
    var err error
    var listener net.Listener
    for port := minPort; port < maxPort; port++ {
        if listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
            break
        }
    }

    if err != nil {
        log.Printf("Can't create p2p host at [%d]\n", port)
        return "", err
    }

    go func() {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Can't Accept for p2p host at [%d]\n", port)
            return
        }

        this.handleConn(conn)
    }()

    return listener.Addr().String(), nil
}

func (this *host) handleConn(conn net.Conn) {
    go func() {
        for {
            request, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                return
            }
            response, err := this.requestHandler.Run(request, conn)
            if err != nil {
                log.Printf("P2P host handle error %s\n", err.Error())
                return
            }
            if _, err = conn.Write(response); err != nil {
                return
            }
        }
    }()
}