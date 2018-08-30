package main

import (
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    metaConfig := dcconf.NewMetaConfig()
    dcutil.SaveJSONConfig(metaConfig)
    dcutil.SaveJSONConfig(dcconf.NewPointConfig(metaConfig))
    dcutil.SaveJSONConfig(dcconf.NewClientConfig(metaConfig))
    dcutil.SaveJSONConfig(dcconf.NewNodeConfig(metaConfig))
    dcutil.SaveJSONConfig(dcconf.NewHTTPCommands(metaConfig))
}