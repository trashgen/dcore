// +build ignore

package tcp

import (
    "net"
    "fmt"
    "log"
    "bufio"
    dcconf "dcore/codebase/modules/config"
)

type TCPModule struct {
    OnShutdown         func()
    OnNewConnection    func(conn net.Conn)
    OnRemoveConnection func(conn net.Conn)

    port           int
    conns          []*net.Conn
    config         *dcconf.TotalConfig
    listener       net.Listener
    newConns       chan net.Conn
    remConns       chan net.Conn
    requests       chan string
    shutdown       chan struct{}
    requestHandler TCPHandler
}

func NewTCPModule(config *dcconf.TotalConfig, requestHandler TCPHandler, maxConnectionNumber int) *TCPModule {
    return &TCPModule{
        conns              : make([]*net.Conn, 0, maxConnectionNumber),
        config             : config,
        newConns           : make(chan net.Conn),
        remConns           : make(chan net.Conn),
        requests           : make(chan string),
        requestHandler     : requestHandler,
        OnShutdown         : func() {},
        OnNewConnection    : func(conn net.Conn){},
        OnRemoveConnection : func(conn net.Conn){}}
}

/**
 * Если 'isRegHost' == true, то параметр 'port' игнорируется - он берется из конфига.
 * TODO : Refactor create Reg Listener
 */
func (this *TCPModule) StartHost(isRegHost bool, port int) {
    var err error
    var listener net.Listener
    if isRegHost {
        for i := 0; i < this.config.Node.MaxAvailableNodesOnMachine; i++ {
            listener, err = net.Listen("tcp", fmt.Sprintf(":%d", this.config.Node.AvailableRegPorts[i]))
            if err != nil {
                continue
            } else {
                break
            }
        }

        if listener == nil {
            log.Fatalf("startHost for Reg : max Nodes on machine is running [%d]\n", this.config.Node.MaxAvailableNodesOnMachine)
        }
    } else {
        listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
        if err != nil {
            log.Fatalf("startHost (Listen) [%d]: %s\n", port, err.Error())
        }
    }

    // TODO : во второй итерации надо научить приостанавливать процесс ожидания
    this.port = port
    this.listener = listener

    go func() {
        this.processConnections()
    }()

    go func() {
        for {
            conn, err := this.listener.Accept()
            // Extra err check here for 'break'
            if err != nil {
                this.handleFatalError(conn,"startHost (Accept)", err)
                break
            }

            this.newConns <- conn
        }
    }()
}

func (this *TCPModule) closeHost() {
    log.Printf("Graceful shutdown host on [%d]\n", this.port)
    this.listener.Close()
    close(this.newConns)
    close(this.remConns)
    close(this.requests)
    close(this.shutdown)
}

func (this *TCPModule) processConnections() {
    for {
        select {
            case newConn := <- this.newConns:
                this.OnNewConnection(newConn)
                this.processMessages(newConn)
            case remConn := <- this.remConns:
                this.OnRemoveConnection(remConn)
            case <- this.shutdown:
                this.OnShutdown()
                this.closeHost()

                return
        }
    }
}

func (this *TCPModule) processMessages(conn net.Conn) {
    response, err := this.requestHandler.Handle(this.getRequest(conn))
    this.handleFatalError(conn, "processMessages", err)
    this.sendResponse(conn, response)
}

func (this *TCPModule) getRequest(conn net.Conn) string {
    out, err := bufio.NewReader(conn).ReadString('\n')
    //this.handleFatalError(conn, "getRequest", err)
    if err != nil {
        this.remConns <- conn
        log.Printf("%s: %s\n", "getRequest", err.Error())
        this.shutdown <- struct{}{}
    }

    return out
}

func (this *TCPModule) sendResponse(conn net.Conn, message string) {
    _, err := conn.Write([]byte(message))
    this.handleFatalError(conn, "sendResponse", err)
}

func (this *TCPModule) handleFatalError(conn net.Conn, methodName string, err error) {
    if err != nil {
        this.remConns <- conn
        log.Printf("%s: %s\n", methodName, err.Error())
        this.shutdown <- struct{}{}
    }
}