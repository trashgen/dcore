package config

import (
    "os"
    "fmt"
    "log"
    "io/ioutil"
    "encoding/json"
)

const MaxAvailableNodesOnMachine = 16

type SignalServerRequestStringConfig struct {
    Check   string
    Remove  string
    ListAll string
}

type SignalServerConfig struct {
    ListenPort int
}

type NodeInfo struct {
    MaxRegPort                 int
    StartRegPort               int
    AvailableRegPorts          [MaxAvailableNodesOnMachine]int
    MaxP2PConnections          int
    RequestActiveNodesCount    int
    MaxAvailableNodesOnMachine int
}

type TotalConfig struct {
    Node            NodeInfo
    SSConfig        SignalServerConfig
    SSCommand       SignalServerRequestStringConfig
    SecretMD5Phrase string
}

func NewTotalConfig() *TotalConfig {
    return &TotalConfig{}
}

func MakeSignalServerConfig() SignalServerConfig {
    return SignalServerConfig{ListenPort:30001}
}

func MakeSignalServerRequestStringConfig() SignalServerRequestStringConfig {
    return SignalServerRequestStringConfig{Remove : "remove", ListAll : "listall", Check : "check"}
}

func MakeNodeInfo() NodeInfo {
    out := NodeInfo{
        StartRegPort               : 57841,
        MaxP2PConnections          : 16,
        RequestActiveNodesCount    : 16,
        MaxAvailableNodesOnMachine : MaxAvailableNodesOnMachine}
    for i := 0; i < out.MaxAvailableNodesOnMachine; i++ {
        out.AvailableRegPorts[i] = out.StartRegPort + i
    }

    out.MaxRegPort = out.AvailableRegPorts[out.MaxAvailableNodesOnMachine - 1]
    return out
}

func (self *TotalConfig) BuildListPort() string {
    return fmt.Sprintf(":%d", self.SSConfig.ListenPort)
}

func (self *TotalConfig) BuildListAllURL() string {
    return fmt.Sprintf("/%s", self.SSCommand.ListAll)
}

func (self *TotalConfig) BuildRemoveURL() string {
    return fmt.Sprintf("/%s", self.SSCommand.Remove)
}

func (self *TotalConfig) BuildCheckURL() string {
    return fmt.Sprintf("/%s", self.SSCommand.Check)
}

func (self *TotalConfig) ReFileWithHardcodedValues() {
    self.SecretMD5Phrase = "operation cwal (C) Starcraft"

    self.Node      = MakeNodeInfo()
    self.SSConfig  = MakeSignalServerConfig()
    self.SSCommand = MakeSignalServerRequestStringConfig()
}

func (self *TotalConfig) SaveConfig() {
    bdata, err := json.MarshalIndent(self, "  ", "\t")
    if err != nil {
        log.Fatal(err.Error())
    }
    
    file, err := os.OpenFile("config/config.cfg", os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0666)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer file.Close()
    
    _, err = file.Write(bdata)
    if err != nil {
        log.Fatal(err.Error())
    }
}

func (self *TotalConfig) LoadConfig() {
    file, err := os.OpenFile("config/config.cfg", os.O_RDONLY,0666)
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