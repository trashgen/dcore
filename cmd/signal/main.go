package main

import (
    "io"
    "log"
    "net/http"
    dcconf "dcore/codebase/config"
)

func main() {
    config := dcconf.NewTotalConfig()
    config.LoadConfig()

    http.HandleFunc(config.BuildListAllURL(), cmdListAll(config))
    http.HandleFunc(config.BuildRegMeURL(), cmdRegMe)
    if err := http.ListenAndServe(config.BuildListPort(), nil); err != nil {
        log.Fatalf("Error starting server:\n%s", err.Error())
    }
}

///////////////////////////////////////////////////////////////////////////////

func cmdListAll(config *dcconf.TotalConfig) (func (w http.ResponseWriter, r *http.Request)) {
    return func (w http.ResponseWriter, r * http.Request) {
        w.Header().Set("Connection", "close")
        w.WriteHeader(http.StatusOK)
        io.WriteString(w, config.BuildConnectedNodesList())
    }
}

func cmdRegMe(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, "RegMe URL requested\n")
}