package main

import (
    "dcore/codebase/modules/p2p"
    "fmt"
    "log"
    "net"
)

type MyRequestHandler struct {}
type MyResponseHandler struct {}

func main() {
    node := p2p.NewNodeModule(&MyRequestHandler{}, &MyResponseHandler{})
    regListenPort := node.Start()
    node.Connect(regListenPort)
    node.Accepting()
}

func (this *MyRequestHandler) Run(data string, conn net.Conn) ([]byte, error) {
    return []byte(fmt.Sprintf("MyRequestHandler: [%s]\n", data)), nil
}

func (this *MyResponseHandler) Run(data string, conn net.Conn) error {
    if _, err := conn.Write([]byte("Fuck off idiot!\n")); err != nil {
        log.Fatalf("MyResponseHandler write: [%s]\n", err.Error())
    }
    log.Printf("MyResponseHandler: [%s]\n", data)
    return nil
}