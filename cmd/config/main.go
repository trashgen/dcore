package main

import (
    dcconf "dcore/codebase/modules/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.ReFileWithHardcodedValues()
    config.SaveConfig()
}
