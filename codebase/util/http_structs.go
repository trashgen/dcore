package util

type ConnectionID struct {
    Key     string
    Address string
}

func NewConnectionID(key string, address string) *ConnectionID {
    return &ConnectionID{Key:key, Address:address}
}

///////////////////////////////////////////////////////////////////////////////

type RequestReg struct {
    Address string
}

type ResponseReg struct {
    Connections map[string]*ConnectionID
    RequestorID string
}

///////////////////////////////////////////////////////////////////////////////

type RequestLook struct {
    Count int
}

type ResponseLook struct {
    Connections map[string]*ConnectionID
}

///////////////////////////////////////////////////////////////////////////////


type RequestPoints struct {
    Count int
}

type ResponsePoints struct {
    Connections map[string]*ConnectionID
}

///////////////////////////////////////////////////////////////////////////////

type RequestRemove struct {
    Key string
}

///////////////////////////////////////////////////////////////////////////////

type RequestCheck struct {
    Key string
}