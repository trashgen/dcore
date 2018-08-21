package main

import (
    "io"
    "log"
    "net/http"
    dcconf "dcore/codebase/config"
)

const (
    ResponseListAll = "Hello, World!"
)

func main() {
    config := dcconf.NewConfig()
    config.LoadConfig()

    http.HandleFunc(buildListAll(config), cmdListAll)
    err := http.ListenAndServe(buildListPort(config), nil)
    if err != nil {
        log.Fatalf("Error starting server:\n%s", err.Error())
    }
}

func cmdListAll(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, ResponseListAll)
}

func buildListPort(config *dcconf.Config) string {
    return ""
}

func buildListAll(config *dcconf.Config) string {
    return ""
}