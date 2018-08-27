package tcp

type TCPHandler interface {
    Handle(message string) (string, error)
}