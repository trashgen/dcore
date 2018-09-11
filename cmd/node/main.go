package main

import "dcore/codebase/module/p2p"

func main() {
    node := p2p.NewNodeModule()
    node.Start()
    node.Accepting()
}