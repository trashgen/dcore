package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dchttp "dcore/codebase/modules/http"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    c, err := dcutil.LoadJSONConfig(dcconf.NewPointConfig(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }

    config, ok := c.(*dcconf.PointConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    c, err = dcutil.LoadJSONConfig(dcconf.NewHTTPCommands(dcconf.NewMetaConfig()))
    if err != nil {
        log.Fatal(err.Error())
    }

    cmdConfig, ok := c.(*dcconf.HTTPCommands)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    dchttp.NewPoint(config, cmdConfig).Start()
}