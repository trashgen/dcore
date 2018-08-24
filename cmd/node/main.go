package main

import (
    dcp2p "dcore/codebase/modules/p2p"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    node := dcp2p.NewNodeModule(config)
    // Эта последовательность методов гарантирует, что хост регистраций успеет подняться.
    go node.StartRegHost()
    node.GetActiveNodeList()

    select{}
}