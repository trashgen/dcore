// +build ignore

package main

import (
    dcp2p "dcore/codebase/modules/p2p"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    node := dcp2p.NewNodeModule(config)
    node.StartBuildingP2P()

    //node.ProcessWork()
    select {}
}