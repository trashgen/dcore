package http

type NodeID struct {
    ID      string
    Port    int
    Address string
}

type ResponseListall struct {
    Nodes       map[string]*NodeID
    RequestorID string
}

type ResponseRemove struct {
    OpResult bool
}

type ResponseCheck struct {
    OpResult bool
}

func NewNodeID(id string, address string, port int) *NodeID {
    return &NodeID{ID:id, Address:address, Port:port}
}

func NewResponseListall(id string, onlineNodes map[string]*NodeID) *ResponseListall {
    return &ResponseListall{RequestorID:id, Nodes:onlineNodes}
}