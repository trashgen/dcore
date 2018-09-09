package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
    dchttpsrv "dcore/codebase/modules/http/server"
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

    dchttpsrv.NewPoint(config, cmdConfig, dchttpsrv.NewMockPersistModule()).Start()
}