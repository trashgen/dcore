package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "strconv"
    dcmisc "dcore/codebase/modules/misc"
    dcview "dcore/codebase/modules/view"
    dcconf "dcore/codebase/modules/config"
    "strings"
)

type NodeModule struct {
    ID             string
    config         *dcconf.TotalConfig
    httpClient     *dcview.ViewModule
    workRegConn    chan net.Conn
    connectedNodes map[string]*dcmisc.NodeID
}

func NewNodeModule(config *dcconf.TotalConfig) *NodeModule {
    return &NodeModule {
        config         : config,
        httpClient     : dcview.NewViewModule(config),
        workRegConn    : make(chan net.Conn),
        connectedNodes : make(map[string]*dcmisc.NodeID, config.Node.RequestActiveNodesCount)}
}

func (this *NodeModule) GetActiveNodeList() {
    data := this.httpClient.GetRawContent(this.config.SSCommand.ListAll, strconv.Itoa(this.config.Node.RequestActiveNodesCount))
    response := this.httpClient.MapListall(data)
    printListall(response)
    this.ID = response.RequestorID
    this.connectedNodes = response.Nodes
}

func printListall(response *dcmisc.RequestListall) {
    sb := strings.Builder{}
    sb.WriteString(fmt.Sprintf("Listall:\n"))
    for _, nodeID := range response.Nodes {
        sb.WriteString(fmt.Sprintf("\tID      = [%s]\n", nodeID.ID))
        sb.WriteString(fmt.Sprintf("\tAddress = [%s]\n", nodeID.Address))
        sb.WriteString(fmt.Sprintf("\tPort    = [%d]\n", nodeID.Port))
        sb.WriteString("\t================================\n")
    }
    log.Print(sb.String())
}

func (this *NodeModule) StartRegHost() {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", this.config.Node.RegisterListenPort))
    if err != nil {
        log.Fatalf("StartRegHost: %s\n", err.Error())
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatalf("StartRegHost: %s\n", err.Error())
        }
    
        go this.handleRegConn(conn)
    }
}

func (this *NodeModule) handleRegConn(conn net.Conn) {
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("handleRegConn: %s\n", err.Error())
    }
    
    response, err := this.processMessage(message)
    if err != nil {
        log.Fatalf("handleRegConn: %s\n", err.Error())
    }
    
    conn.Write([]byte(response))
}

func (this *NodeModule) processMessage(message string) (string, error) {
    return "", nil
}