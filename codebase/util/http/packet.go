package http

import "fmt"

type ConnectionID struct {
    IP   string
    Key  string
    Port int
}

func NewConnectionID(key string, ip string, port int) *ConnectionID {
    return &ConnectionID{Key:key, IP:ip, Port:port}
}

func (this ConnectionID) Address() string {
    return fmt.Sprintf("%s:%d", this.IP, this.Port)
}

///////////////////////////////////////////////////////////////////////////////

type RequestReg struct {
    Port int
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