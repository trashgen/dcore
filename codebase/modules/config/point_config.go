package config

import "fmt"

type CommandDesc struct {
    Name  string
    Param string
}

type PointConfig struct {
    Reg          CommandDesc
    Look         CommandDesc
    Root         CommandDesc
    Check        CommandDesc
    Points       CommandDesc
    Remove       CommandDesc
    FileName     string
    ListenPort   int
    SecretPhrase string
}

func NewPointConfig(fileName string) *PointConfig {
    return &PointConfig{
        Reg          : CommandDesc{Name:"reg"},
        Look         : CommandDesc{Name:"look", Param:"count"},
        Root         : CommandDesc{Name:""},
        Check        : CommandDesc{Name:"check", Param:"key"},
        Points       : CommandDesc{Name:"points", Param:"count"},
        Remove       : CommandDesc{Name:"remove", Param:"key"},
        FileName     : fileName,
        ListenPort   : 30001,
        SecretPhrase : "operation cwal (C) Starcraft"}
}

func (this *PointConfig) FormattedListenPort() string {
    return fmt.Sprintf(":%d", this.ListenPort)
}