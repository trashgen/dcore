package server

import (
    "fmt"
    "strconv"
)

func RegPacket1013ID() int   {return 1013}
func DeathPacket777ID() int  {return 777}
func ConfirmPacket88ID() int {return 88}

// TODO : расширить адресом поинта
func BuildRequest1013(key string) []byte {
    return []byte(fmt.Sprintf("1013\t%s\n", key))
}

// TODO : расширить адресом поинта
func BuildGoodResponse1013(status bool, key string, address string) []byte {
    return []byte(fmt.Sprintf("1013\t%s\t%s\t%s\n", strconv.FormatBool(status), key, address))
}

func BuildBadResponse1013(status bool) []byte {
    return []byte(fmt.Sprintf("1013\t%s\n", strconv.FormatBool(status)))
}

func BuildCommand777(status bool) []byte {
    return []byte(fmt.Sprintf("777\t%s\n", strconv.FormatBool(status)))
}

func BuildRequest88(key string, addr string) []byte {
    return []byte(fmt.Sprintf("88\t%s\t%s\n", key, addr))
}

func BuildResponse88(key string) []byte {
    return []byte(fmt.Sprintf("88\t%s\n", key))
}