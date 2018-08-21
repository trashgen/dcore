package main

import (
    dcconf "dcore/codebase/config"
)

func main() {
    config := dcconf.NewConfig()
    config.ReFileWithHardcodedValues()
}
