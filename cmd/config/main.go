package main

import (
    dcutil "dcore/codebase/util"
    dcconf "dcore/codebase/modules/config"
)

func main() {
    // TODO : 'configFileName' to metaconfig
    configFileName := "pointconfig.cfg"
    dcutil.SaveJSONConfig(configFileName, dcconf.NewPointConfig(configFileName))

    configFileName = "clientconfig.cfg"
    dcutil.SaveJSONConfig(configFileName, dcconf.NewClientConfig(configFileName))
}
