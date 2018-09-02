package p2p

import (
    "fmt"
    "log"
    "net"
    "time"
    "bufio"
    "dcore/codebase/modules/config"
    dctcpsrvutil "dcore/codebase/util/tcp/server"
)

type Host struct {}

func NewP2PHost() *Host {
    return &Host{}
}

// TODO : Ограничить на одно соединение
func (this *Host) Start(config *config.NodeConfig) (string, bool) {
    // TODO : создает хост и отдает полный адрес для подключения
    var port int
    var err error
    var listener net.Listener
    for port := config.MinP2PPort; port < config.MaxP2PPort; port++ {
        if listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
            break
        }
    }

    if err != nil {
        log.Printf("Can't create p2p host at [%d]\n", port)
        return "", false
    }

    go func() {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Can't Accept for p2p host at [%d]\n", port)
            return
        }

        this.handleConn(conn)
    }()

    log.Printf("Host address : [%s]\n", listener.Addr().String())
    return listener.Addr().String(), true
}

func (this *Host) handleConn(conn net.Conn) {
    go func() {
        for {
            response, err := bufio.NewReader(conn).ReadString('\n')
            if err != nil {
                log.Printf("P2P Client receive error %s\n", err.Error())
                return
            }
            request := dctcpsrvutil.BuildPacket111(conn.RemoteAddr().String())
            conn.Write(request)
            time.Sleep(time.Second)
        }
    }()
}