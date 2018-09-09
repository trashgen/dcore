package http

import "reflect"

const (
    RegName = "reg"
)

type Reg struct {
    Name     string
    Port     int
    PortName string
    PortType reflect.Kind
}
func NewReg() *Reg {
    return &Reg{Name: RegName, PortName: "port", PortType: reflect.Int}
}