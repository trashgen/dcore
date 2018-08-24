package misc

type NodeID struct {
    ID      string
    Port    int
    Address string
}

type RequestListall struct {
    Nodes       map[string]*NodeID
    RequestorID string
}

type RequestRemove struct {
    OpResult bool
}

type RequestCheck struct {
    OpResult bool
}

func NewNodeID(id string, address string, port int) *NodeID {
    return &NodeID{ID:id, Address:address, Port:port}
}