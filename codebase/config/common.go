package config

import (
    "os"
    "log"
    "path"
    "encoding/json"
    "path/filepath"
    "io/ioutil"
    "strings"
    "fmt"
)

type CommonConfig struct {
    SSListenPort int
    SSListAll    string
}

type ViewConfig struct {
    SSHost string
}

type Config struct {
    Common CommonConfig
    View   ViewConfig
}

func NewConfig() *Config {
    return &Config{}
}

func (self *Config) ReFileWithHardcodedValues() {
    config := Config{
        View: ViewConfig{SSHost:"127.0.0.1"},
        Common: CommonConfig{
            SSListenPort:30001,
            SSListAll:"listall"}}

    bdata, err := json.MarshalIndent(config, "  ", "  ")
    //bdata, err := json.Marshal(config)
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

func (self *Config) Dump() {
    sb := strings.Builder{}
    sb.WriteString("===============================\n")
    sb.WriteString(fmt.Sprintf("Common.SSListenPort : [%d]\n", self.Common.SSListenPort))
    sb.WriteString(fmt.Sprintf("Common.SSListAll    : [%s]\n", self.Common.SSListAll))
    sb.WriteString(fmt.Sprintf("View.SSHost         : [%s]\n", self.View.SSHost))
    //sb.WriteString(fmt.Sprintf("", ))
    sb.WriteString("===============================\n")

    fmt.Printf(sb.String())
}