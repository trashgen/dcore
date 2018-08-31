package server

import (
    "fmt"
    "time"
    "strings"
    "unicode"
    "crypto/md5"
    "encoding/hex"
    "dcore/codebase/modules/config"
    dchttputil "dcore/codebase/util/http"
)

func BuildRegResponse(remoteAddr string, secret string) (response []byte, key string, ip string) {
    response = make([]byte, 0)

    hash := md5.New()
    hash.Write([]byte(fmt.Sprintf("%s!@#$^&*()%s!@#$^&*()%s", remoteAddr, time.Now().String(), secret)))
    key = hex.EncodeToString(hash.Sum(nil))
    ip = strings.TrimSuffix(strings.TrimRightFunc(remoteAddr,
        func(r rune) bool {
            return unicode.IsDigit(r)
        }), ":")

    return []byte(key), key, ip
}

func BuildLookOrPointsResponse(nodes map[string]*dchttputil.ConnectionID, count int) []byte {
    if count > len(nodes) {
        count = len(nodes)
    }

    i := 0
    sb := strings.Builder{}
    for _, v := range nodes {
        sb.WriteString(fmt.Sprintf("%s-%s:%d\t", v.Key, v.IP, v.Port))
        if i++; i == count {
            break
        }
    }

    return []byte(strings.TrimSuffix(sb.String(), "\t"))
}

func BuildRootResponse(cmdConfig *config.HTTPCommands) []byte {
    sb := strings.Builder{}
    sb.WriteString("<h1>Point help:</h1>")
    sb.WriteString("<b>Root page</b>: You are here now<br>")
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to register on Point. No query params. Key (string) as result<br>", cmdConfig.Reg.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Nodes. If query param 'count' here with (int) > 0 as value - then limit number of Nodes to send in Response<br>", cmdConfig.Look.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to check if Node is registered at this Point. Key (string) as query param required<br>", cmdConfig.Check.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request list of active Points. If query param 'count' here with (int) > 0 as value - then limit number of Points to send in Response<br>", cmdConfig.Points.Name))
    sb.WriteString(fmt.Sprintf("<b>%s</b>: Request to remove Node. Key (string) as query param required<br>", cmdConfig.Remove.Name))

    return []byte(sb.String())
}

func BuildCheckOrRemoveResponse(nodes map[string]*dchttputil.ConnectionID, key string) (response []byte, ok bool) {
    var result string
    _, ok = nodes[key]
    if ok {
        result = "true"
    } else {
        result = "false"
    }

    return []byte(result), ok
}