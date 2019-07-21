package p2p

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"

	dcutil "dcore/codebase/util"
	dcutcp "dcore/codebase/util/tcp/server"
)

type regHostModule struct {
	*mediator
	newConn     chan net.Conn
	regListener net.Listener
}

func newRegHostModule(m *mediator) *regHostModule {
	return &regHostModule{
		newConn:  make(chan net.Conn),
		mediator: m}
}

func (this *regHostModule) startRegHost() (port int) {
	var err error
	for port = this.nodeConfig.MinRegPort; port < this.nodeConfig.MaxRegPort; port++ {
		if this.regListener, err = net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("Can't start reg host [%s]\n", err.Error())
	}
	this.onNewConnection()
	return port
}

func (this *regHostModule) Accepting() {
	for {
		conn, err := this.regListener.Accept()
		if err != nil {
			log.Fatalf("Can't Accept new connections: [%s]\n", err.Error())
		}
		this.newConn <- conn
	}
}

func (this *regHostModule) onNewConnection() {
	go func() {
		for conn := range this.newConn {
			func(c net.Conn) {
				this.createP2PLine(c)
			}(conn)
		}
	}()
}

func (this *regHostModule) createP2PLine(conn net.Conn) {
	go func() {
		var err error
		for {
			var data string
			data, err = bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				break
			}
			var response []byte
			var hasResponseData bool
			response, hasResponseData, err = this.handle(data, conn)
			if err != nil {
				log.Printf("Reg host handler error [%s]: [%s]\n", data, err.Error())
				break
			}
			if hasResponseData {
				if _, err = conn.Write(response); err != nil {
					break
				}
			} else {
				break
			}
		}
	}()
}

func (this *regHostModule) handle(data string, conn net.Conn) ([]byte, bool, error) {
	packetID, params, err := dcutil.SplitPacketIDWithData(data)
	if err != nil {
		return nil, false, err
	}
	switch packetID {
	case dcutcp.RegPacket1013ID:
		return this.handle1013Request(params, conn.RemoteAddr())
	case dcutcp.DeathPacket777ID:
		return this.handle777Command(params)
	case dcutcp.ConfirmPacket88ID:
		return this.handle88Command(params)
	}
	return nil, false, nil
}

func (this *regHostModule) handle1013Request(params []string, address net.Addr) ([]byte, bool, error) {
	var err error
	request, err := dcutil.Split1013RequestParams(params)
	if err != nil {
		return nil, false, err
	}

	var response []byte
	status := this.clientModule.RequestCheck(this.ThisNodeKey, request.Target)
	if status {
		l := newLine(this.mediator, request.Target)
		l.startHost()
		response = dcutcp.BuildGoodResponse1013(status, this.ThisNodeKey, l.address)
	} else {
		response = dcutcp.BuildBadResponse1013(status)
		err = errors.New("can't reg node with invalid key")
		this.clientModule.SendRequestBan(this.ThisNodeKey, dcutil.RemovePortFromAddressString(address.String()))
	}

	return response, true, err
}

func (this *regHostModule) handle777Command(params []string) ([]byte, bool, error) {
	if command, err := dcutil.SplitCommand777Params(params); !command.Status || err != nil {
		log.Fatal("I am LIAR!!!!")
	}
	return nil, false, nil
}

func (this *regHostModule) handle88Command(params []string) ([]byte, bool, error) {
	if command, err := dcutil.SplitCommand88Params(params); err == nil {
		l := this.lines[command.ThoseNodeKey]
		l.startClient(command.HostAddr)
	}
	return nil, false, nil
}
