package main

import (
    "log"
    dcutil "dcore/codebase/util"
    dcp2p "dcore/codebase/modules/p2p"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config, ok := dcutil.LoadJSONConfig(dcconf.NewNodeConfig(dcconf.NewMetaConfig())).(*dcconf.NodeConfig)
    if ! ok {
        log.Fatal("Config: type mismatch")
    }

    node := dcp2p.NewNodeModule(config)
    node.Connect()
    node.Start()
}