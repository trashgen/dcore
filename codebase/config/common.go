package config

import (
    "os"
    "log"
    "path"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
)

type SignalServerRequestStringConfig struct {
    ListAll string
    RegMe   string
}

type SignalServerConfig struct {
    ListenPort int
}

type AddressInfo struct {
    IP   string
    Port int
}

type Config struct {
    Nodes      []*AddressInfo
    Signals    []*AddressInfo
    SSConfig      SignalServerConfig
    SSCommand     SignalServerRequestStringConfig
}

func NewConfig() *Config {
    return &Config{Nodes : make([]*AddressInfo, 0, 2)}
}

func NewAddressInfo(ip string, port int) *AddressInfo {
    return &AddressInfo { IP:ip, Port:port}
}

func (self *Config) ReFileWithHardcodedValues() {
    self.SSConfig  = SignalServerConfig{ListenPort:30001}
    self.SSCommand = SignalServerRequestStringConfig{RegMe : "regme", ListAll : "listall"}

    self.Nodes      = make([]*AddressInfo, 0, 2)
    self.Nodes      = append(self.Nodes, NewAddressInfo("127.0.0.1", 30001))
    self.Nodes      = append(self.Nodes, NewAddressInfo("I am bad IP", 666))

    self.Signals    = make([]*AddressInfo, 0, 2)
    self.Signals    = append(self.Signals, NewAddressInfo("127.0.0.1", 30001))
    self.Signals    = append(self.Signals, NewAddressInfo("I am bad IP", 666))
    
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