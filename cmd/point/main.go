package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config, ok := dcutil.LoadJSONConfig(dcconf.NewPointConfig(dcconf.NewMetaConfig())).(*dcconf.PointConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    cmdConfig, ok := dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig())).(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    dchttp.NewPoint(config, cmdConfig).Start()
}