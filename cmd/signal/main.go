package main

import (
    "io"
    "log"
    "net/http"
)

const (
    AddressListAll  = "/listall"
    ResponseListAll = "Hello, World!"
)

func main() {
    http.HandleFunc(AddressListAll, cmdListAll)
    err := http.ListenAndServe(":30001", nil)
    if err != nil {
        log.Fatalf("Error starting server:\n%s", err.Error())
    }
}

func cmdListAll(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Connection", "close")
    w.WriteHeader(http.StatusOK)
    io.WriteString(w, ResponseListAll)
}