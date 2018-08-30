package config

type ClientConfig struct {
    Reg         CommandDesc
    Look        CommandDesc
    Root        CommandDesc
    Check       CommandDesc
    Points      CommandDesc
    Remove      CommandDesc
    FileName    string
    EntryPoints []string
}

func NewClientConfig(fileName string) *ClientConfig {
    points := make([]string, 0)
    points = append(points, "http://localhost:30001")
    points = append(points, "КРОВЬКИШКИРАСПИДОРАСИЛО:11111")

    return &ClientConfig{
        Reg         : CommandDesc{Name:"reg"},
        Look        : CommandDesc{Name:"look", Param:"count"},
        Root        : CommandDesc{Name:""},
        Check       : CommandDesc{Name:"check", Param:"key"},
        Points      : CommandDesc{Name:"points", Param:"count"},
        Remove      : CommandDesc{Name:"remove", Param:"key"},
        FileName    : fileName,
        EntryPoints : points}
}