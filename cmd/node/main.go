package main

import "dcore/codebase/modules/p2p"

func main() {
    node := p2p.NewNodeModule()
    node.Start()
    node.Connect()
    node.Accepting()
}