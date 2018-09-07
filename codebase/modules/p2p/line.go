package p2p

import (
    "log"
    "fmt"
)

// Представляет собой полное описание одного P2P соединения с обработкой всех команд
type line struct {
    *host
    *client
    // *handler
    *mediator
    // TODO : В итоге будет JSON
    inCommand  chan string
    outCommand chan string
    // Struct unique key
    thoseNodeKey string
}

func newLine(m *mediator, thoseNodeKey string) *line {
    out := &line{
        mediator     : m,
        inCommand    : make(chan string),
        outCommand   : make(chan string),
        thoseNodeKey : thoseNodeKey}
    out.lines[thoseNodeKey] = out
    return out
}

func (this *line) Send(targetNodeKey string) {
    l := this.lines[targetNodeKey]
    l.inCommand <- "Hello, World!\n"
}

func (this *line) startHost() {
    this.host = newHost(this.nodeConfig.MinP2PPort, this.nodeConfig.MaxP2PPort, this.inCommand)
    if err := this.host.Start(); err != nil {
        log.Fatal(err.Error())
    }
    log.Printf("Start Host on [%d] with key [%s]\n", this.host.port, this.thoseNodeKey)
}

func (this *line) startClient(thoseHostAddr string) {
    this.client = newClient(this.outCommand)
    if err := this.client.Connect(thoseHostAddr); err != nil {
        // TODO : обработать невозможность прямого подключения
        log.Fatal(err.Error())
    }
    log.Printf("Start Client on [%s] with key [%s]\n", thoseHostAddr, this.thoseNodeKey)
    this.handleCommand()
}

func (this *line) handleCommand() {
    go func() {
        for inCommand := range this.inCommand {
            response, hasResponse := doSomeIn(inCommand)
            if hasResponse {
                this.outCommand <- response
            }
        }
    }()
}

func doSomeIn(inCommand string) (string, bool) {
    return fmt.Sprintf("Response for: [%s]\n", inCommand), true
}