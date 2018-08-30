package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    // TODO : 'configFileName' to metaconfig
    configFileName := "clientconfig.cfg"
    config, ok := dcutil.LoadJSONConfig(configFileName, dcconf.NewClientConfig(configFileName)).(*dcconf.ClientConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    httpClient := dchttp.NewClientModule(config)
    //data := httpClient.RequestLook(2, 3)
    data := httpClient.RequestReg()
    log.Printf("reg Response is\n[%s]\n", data)
}