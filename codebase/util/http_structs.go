package util

type ConnectionID struct {
    IP   string
    Key  string
    Port int
}

type RequestReg struct {
    Count int
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