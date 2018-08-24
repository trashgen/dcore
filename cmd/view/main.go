package main

import (
    "fmt"
    "log"
    "flag"
    "strconv"
    "strings"
    dcmisc "dcore/codebase/modules/misc"
    dcview "dcore/codebase/modules/view"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    method, param := getCmdParams(config)

    httpClient := dcview.NewViewModule(config)
    data := httpClient.GetRawContent(method, param)

    switch method {
        case config.SSCommand.ListAll:
            printListall(httpClient.MapListall(data))
        case config.SSCommand.Remove:
            printRemove(httpClient.MapRemove(data))
        case config.SSCommand.Check:
            printCheck(httpClient.MapCheck(data))
    }
}

func getCmdParams(config *dcconf.TotalConfig) (string, string) {
    method := flag.String("method", config.SSCommand.ListAll, "only [listall|remove|check]")
    queryParam := flag.String("param", "0", "respectivly [int|string|string]")
    flag.Parse()

    return *method, *queryParam
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

func printRemove(response *dcmisc.RequestRemove) {
    sb := strings.Builder{}
    sb.WriteString(fmt.Sprintf("Remove:\n"))
    sb.WriteString(fmt.Sprintf("\tOpResult = [%s]\n", strconv.FormatBool(response.OpResult)))
    log.Print(sb.String())
}

func printCheck(response *dcmisc.RequestCheck) {
    sb := strings.Builder{}
    sb.WriteString(fmt.Sprintf("Check:\n"))
    sb.WriteString(fmt.Sprintf("\tOpResult = [%s]\n", strconv.FormatBool(response.OpResult)))
    log.Print(sb.String())
}
