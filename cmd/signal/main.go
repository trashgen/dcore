package main

import (
    "log"
    "net/http"
    dcconf "dcore/codebase/modules/config"
    dchttp "dcore/codebase/modules/http"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    if err := http.ListenAndServe(config.BuildListPort(), dchttp.NewHTTPModule(config)); err != nil {
        log.Fatalf("Error starting server: %s", err.Error())
    }
}