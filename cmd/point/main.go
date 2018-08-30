package main

import (
    "log"
    "net/http"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    // TODO : 'configFileName' to metaconfig
    configFileName := "pointconfig.cfg"
    config, ok := dcutil.LoadJSONConfig(configFileName, dcconf.NewPointConfig(configFileName)).(*dcconf.PointConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    if err := http.ListenAndServe(config.FormattedListenPort(), dchttp.NewPoint(config)); err != nil {
        log.Fatalf("Error starting server: %s", err.Error())
    }
}