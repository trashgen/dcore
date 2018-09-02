package server

import (
    "fmt"
    "strconv"
)

func RegPacketID() int {return 1013}
func DeathPacketID() int {return 777}
func ConfirmPacketID() int {return 88}

// TODO : расширить адресом поинта
func BuildPacket1013Request(key string) []byte {
    return []byte(fmt.Sprintf("1013\t%s\n", key))
}

// TODO : расширить адресом поинта
func BuildGoodPacket1013Response(status bool, key string, address string) []byte {
    return []byte(fmt.Sprintf("1013\t%s\t%s\t%s\n", strconv.FormatBool(status), key, address))
}

func BuildBadPacket1013Response(status bool) []byte {
    return []byte(fmt.Sprintf("1013\t%s\n", strconv.FormatBool(status)))
}

func BuildPacket777(status bool) []byte {
    return []byte(fmt.Sprintf("777\t%s\n", strconv.FormatBool(status)))
}

func BuildPacket88(addr string) []byte {
    return []byte(fmt.Sprintf("88\t%s\n", addr))
}

func BuildPacket111(addr string) []byte {
    return []byte(fmt.Sprintf("111\tHello, World to %s!!\n", addr))
}