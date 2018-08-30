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

    dchttp.NewPoint(config).Start()
}