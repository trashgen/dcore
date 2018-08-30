package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "strings"
    dcmisc "dcore/codebase/modules/misc"
    dcconf "dcore/codebase/modules/config"
)

type RegClientModule struct {
    config        *dcconf.NodeConfig
    otherRegHosts []*NodeDesc
}

func NewRegClientModule(config *dcconf.NodeConfig) *RegClientModule {
    return &RegClientModule{config:config, otherRegHosts:make([]*NodeDesc, 0, config.MaxP2PConnections)}
}

// RegClient struct
func (this *RegClientModule) Connect() {
    for _, nodeDesc := range this.otherRegHosts {
        go func(nd *NodeDesc) {
            this.connectReghost(nd.Address)
        }(nodeDesc)
    }
}

// RegClient struct
func (this *RegClientModule) parseLookRequest(data string) {
    data = strings.TrimSuffix(data, "\n")
    scanner := bufio.NewScanner(strings.NewReader(data))
    scanner.Split(dcmisc.SplitterFunc('\t'))
    for scanner.Scan() {
        this.otherRegHosts = append(this.otherRegHosts, NewNodeDesc(scanner.Text()))
    }
}

// RegClient struct
func (this *RegClientModule) connectReghost(address string) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", address)
    }
    
    request := this.createPacket1013()
    
    _, err = conn.Write(request)
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }
    
    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }
    
    log.Printf("I am Response from reg host [%s]\n", data)
}

// RegClient struct
func (this *RegClientModule) createPacket1013() []byte {
    return []byte(fmt.Sprintf("Hi! I am Packet 1013!!!"))
}