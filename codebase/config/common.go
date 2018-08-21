package config

import (
    "os"
    "log"
    "path"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
)

type SignalServerConfig struct {
    ListenPort int
}

type AddressInfo struct {
    IP   string
    Port int
}

type Config struct {
    SSConfig     SignalServerConfig
    Nodes      []*AddressInfo
    Signals    []*AddressInfo
    SSCommands []string
}

func NewConfig() *Config {
    return &Config{Nodes : make([]*AddressInfo, 0, 100)}
}

func NewAddressInfo(ip string, port int) *AddressInfo {
    return &AddressInfo { IP:ip, Port:port}
}

func (self *Config) ReFileWithHardcodedValues() {
    self.SSConfig = SignalServerConfig{ListenPort:30001}
    
    self.Nodes      = make([]*AddressInfo, 0, 2)
    self.Signals    = make([]*AddressInfo, 0, 2)
    self.SSCommands = make([]string, 0, 2)

    self.Nodes      = append(self.Nodes, NewAddressInfo("127.0.0.1", 30001))
    self.Nodes      = append(self.Nodes, NewAddressInfo("I am bad IP", 666))
    self.Signals    = append(self.Signals, NewAddressInfo("127.0.0.1", 30001))
    self.Signals    = append(self.Signals, NewAddressInfo("I am bad IP", 666))
    self.SSCommands = append(self.SSCommands, "listall")
    self.SSCommands = append(self.SSCommands, "regme")

    bdata, err := json.MarshalIndent(self, "  ", "\t")
    if err != nil {
        log.Fatal(err.Error())
    }
    absPath, err := filepath.Abs("../../bin/dcore/config")
    if err != nil {
        log.Fatal(err.Error())
    }

    file, err := os.OpenFile(path.Join(absPath, "config.cfg"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0666)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer file.Close()

    _, err = file.Write(bdata)
    if err != nil {
        log.Fatal(err.Error())
    }
}

func (self *Config) LoadConfig() {
    absPath, err := filepath.Abs("../../bin/dcore/config")
    if err != nil {
        log.Fatal(err.Error())
    }
    
    file, err := os.OpenFile(path.Join(absPath, "config.cfg"), os.O_RDONLY,0666)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err.Error())
    }

    if err := json.Unmarshal(bytes, self); err != nil {
        log.Fatal(err.Error())
    }
}