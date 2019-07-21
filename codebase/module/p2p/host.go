package p2p

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type host struct {
	port      int
	address   string
	maxPort   int
	minPort   int
	listener  net.Listener
	inCommand chan<- string
}

func newHost(minPort int, maxPort int, in chan<- string) *host {
	return &host{minPort: minPort, maxPort: maxPort, inCommand: in}
}

func (this *host) Start() error {
	var err error
	for this.port = this.minPort; this.port < this.maxPort; this.port++ {
		if this.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", this.port)); err == nil {
			break
		}
	}

	if err != nil {
		return errors.New(fmt.Sprintf("error create p2p host [%s]\n", err.Error()))
	}

	this.address = this.listener.Addr().String()
	go func() {
		if conn, err := this.listener.Accept(); err != nil {
			return
		} else {
			this.handleConn(conn)
		}
	}()
	return nil
}

func (this *host) handleConn(conn net.Conn) {
	go func() {
		for {
			request, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				break
			}
			this.inCommand <- request
		}
		this.close()
	}()
}

func (this *host) close() {
	this.listener.Close()
}
