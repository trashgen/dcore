package p2p

import (
    "log"
    "net"
    "time"
    "bufio"
    dctcpsrvutil "dcore/codebase/util/tcp/server"
)

type Client struct {}

func NewP2PClient() *Client {
    return &Client{}
}

func (this *Client) Connect(address string) bool {
    hostConn, err := net.Dial("tcp", address)
    if err != nil {
        log.Printf("Can't connect to p2p host [%s]\n", address)
        return false
    }

    go func() {
        for {
            request := dctcpsrvutil.BuildPacket111(address)
            hostConn.Write(request)
            response, err := bufio.NewReader(hostConn).ReadString('\n')
            if err != nil {
                log.Printf("P2P Client receive error %s\n", err.Error())
                return
            }
            log.Printf("Response on Client: [%s]\n", response)
            time.Sleep(time.Second)
        }
    }()

    return true
}