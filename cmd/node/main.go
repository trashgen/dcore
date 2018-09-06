package main

import (
    "dcore/codebase/modules/p2p"
)

func main() {
    node := p2p.NewNodeModule()
    regListenPort := node.Start()
    node.Connect(regListenPort)
    node.Accepting()
}