package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config, ok := dcutil.LoadJSONConfig(dcconf.NewClientConfig(dcconf.NewMetaConfig())).(*dcconf.ClientConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }
    
    cmdConfig, ok := dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig())).(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    httpClient := dchttp.NewClientModule(config, cmdConfig)
    data := httpClient.RequestRoot()
    log.Print("================ TESTING HTTP ================\n")
    log.Printf("Response Root is\n[%s]\n", data)
    data = httpClient.RequestLook(1, 3)
    log.Printf("Response Look is\n[%s]\n", data)
    data = httpClient.RequestReg("tratata")
    log.Printf("Response Reg is\n[%s]\n", data)
    data = httpClient.RequestCheck("nodeKey1")
    log.Printf("Response Check is\n[%s]\n", data)
    data = httpClient.RequestRemove("nodeKey1")
    log.Printf("Response Remove is\n[%s]\n", data)
    data = httpClient.RequestCheck("nodeKey1")
    log.Printf("Response Check is\n[%s]\n", data)
    data = httpClient.RequestLook(1, 3)
    log.Printf("Response Look is\n[%s]\n", data)
    data = httpClient.RequestPoints(1, 3)
    log.Printf("Response Points is\n[%s]\n", data)
}