// +build ignore

package p2p

import (
    "fmt"
    "log"
    "net"
    "strconv"
    "strings"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

type RegModule struct {
    ID       string
    config   *dcconf.TotalConfig
    ssClient *dchttp.SSClient
}

func NewRegModule(config *dcconf.TotalConfig, ssClient *dchttp.SSClient) *RegModule {
    return &RegModule{config:config, ssClient:ssClient}
}

func (this *RegModule) Execute() {
    response := this.getActiveNodeList()
    this.ID = response.RequestorID
    for _, v := range response.Nodes {
        address := v.Address
        go this.registerOn(address)
    }
}

func (this *RegModule) registerOn(ip string) {
    var err error
    var regConn net.Conn
    for i := 0; i < this.config.Node.MaxAvailableNodesOnMachine; i++ {
        regConn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", ip, this.config.Node.AvailableRegPorts[i]))
        if err != nil {
            // TODO : Здесь будет добавление в массив подключений через Proxy. Пока под вопросом. Потом проработать алгоритм.
            continue
        } else {
            break
        }
    }

    if regConn == nil {
        log.Fatalf("StartP2PRegistrationProcess (can't dial to reg host %s): %s\n", ip, err.Error())
    }

    // Коннект нужен только на процесс регистрации
    defer regConn.Close()
}

func (this *RegModule) getActiveNodeList() *dchttp.ResponseListall {
    data := this.ssClient.GetRawContent(this.config.SSCommand.ListAll, strconv.Itoa(this.config.Node.RequestActiveNodesCount))
    response := this.ssClient.MapListall(data)
    // TODO : just for testing. Remove It before PUSH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
    printListall(response)
    return response
}

func printListall(response *dchttp.ResponseListall) {
    sb := strings.Builder{}
    sb.WriteString(fmt.Sprintf("Listall:\n"))
    sb.WriteString(fmt.Sprintf("\tKey     = [%s]\n", response.RequestorID))
    for _, nodeID := range response.Nodes {
        sb.WriteString(fmt.Sprintf("\tID      = [%s]\n", nodeID.ID))
        sb.WriteString(fmt.Sprintf("\tAddress = [%s]\n", nodeID.Address))
        sb.WriteString(fmt.Sprintf("\tPort    = [%d]\n", nodeID.Port))
        sb.WriteString("\t================================\n")
    }
    log.Print(sb.String())
}
