package p2p

import (
    "fmt"
    "log"
    "net"
    "bufio"
    "strings"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
)

type regClientModule struct {
    Key           string
    config        *dcconf.NodeConfig
    otherRegHosts []*nodeDesc
}

func newRegClientModule(config *dcconf.NodeConfig) *regClientModule {
    return &regClientModule{config:config, otherRegHosts:make([]*nodeDesc, 0, config.MaxP2PConnections)}
}

func (this *regClientModule) Connect() {
    for _, nd := range this.otherRegHosts {
        go func(nd *nodeDesc) {
            // TODO : проверить - шарица поле Адрес или только сам стракт
            this.connectReghost(nd.Address)
        }(nd)
    }
}

func (this *regClientModule) parseLookResponse(data string) {
    data = strings.TrimSuffix(data, "\n")
    values := dcutil.ScanString(data, '\t')
    for _, nd := range values {
        this.otherRegHosts = append(this.otherRegHosts, newNodeDesc(nd))
    }
}

func (this *regClientModule) connectReghost(address string) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        log.Fatalf("Can't connect to reg host [%s]\n", address)
    }

    _, err = conn.Write(this.createPacket1013Request())
    if err != nil {
        log.Fatalf("Error send Packet 1013 to reg host [%s]\n", err.Error())
    }

    data, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
        log.Fatalf("Error receive Packet 1013 from reg host [%s]\n", err.Error())
    }

    data = strings.TrimSuffix(data, "\n")
    if data == "false" {
        log.Fatal("I am bad boy and my IP in blacklist now :(\n")
    } else {
        log.Printf("I am Response from reg host [%s]\n", data)
    }
}

func (this *regClientModule) createPacket1013Request() []byte {
    return []byte(fmt.Sprintf("1013\t%s\n", this.Key))
}