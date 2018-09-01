package server

import "fmt"

func BuildGoodPacket1013Response(address string) string {
    return fmt.Sprintf("true\t%s\n", address)
}