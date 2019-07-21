package main

import (
	dcconf "dcore/codebase/module/config"
	dcutil "dcore/codebase/util"
)

func main() {
	metaConfig := dcconf.NewMetaConfig()
	dcutil.SaveJSONConfig(metaConfig)
	dcutil.SaveJSONConfig(dcconf.NewPointConfig(metaConfig))
	dcutil.SaveJSONConfig(dcconf.NewClientConfig(metaConfig))
	dcutil.SaveJSONConfig(dcconf.NewNodeConfig(metaConfig))
}
