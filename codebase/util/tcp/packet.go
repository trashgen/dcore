package tcp

type Request1013 struct {
    ID           int
    ThoseNodeKey string
}

type Response1013 struct {
    ID           int
    Status       bool
    Address      string
    ThoseNodeKey string
}

///////////////////////////////////////////////////////////////////////////////

type Request777 struct {
    ID     int
    Status bool
}

type Request88 struct {
    ID           int
    HostAddr     string
    ThoseNodeKey string
}

type RequestResponse111 struct {
    ID           int
    Message      string
    ThoseNodeKey string
}