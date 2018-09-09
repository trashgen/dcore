package util

import (
    "fmt"
    "bufio"
    "bytes"
    "errors"
    "strconv"
    "strings"
    "unicode"
    dcutcp "dcore/codebase/util/tcp"
)

func SplitterFunc(splitSymbol byte) (func(data []byte, atEOF bool) (advance int, token []byte, err error)) {
    return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
        if atEOF && len(data) == 0 {
            return 0, nil, nil
        }
        if i := bytes.IndexByte(data, splitSymbol); i >= 0 {
            return i + 1, data[0:i], nil
        }
        if atEOF {
            return len(data), data, nil
        }
        return 0, nil, nil
    }
}

func SplitPacketIDWithData(data string) (int, []string, error) {
    params := SplitTCPParams(strings.TrimSuffix(data, "\n"))
    if len(params) < 1 {
        return 0, nil, errors.New(fmt.Sprintf("no packet id in request [%s]", data))
    }
    packetID, err := strconv.Atoi(params[0])
    if err != nil {
        return 0, nil, err
    }
    return packetID, params[1:], nil
}

func Split1013RequestParams(params []string) (*dcutcp.Request1013, error) {
    if len(params) != 1 {
        return nil, errors.New(fmt.Sprintf("bad 1013 request [%#v]", params))
    }
    return &dcutcp.Request1013{ID:1013, ThoseNodeKey:params[0]}, nil
}

func Split1013Response(params []string) (*dcutcp.Response1013, error) {
    if len(params) != 3 {
        return nil, errors.New(fmt.Sprintf("bad 1013 response [%#v]", params))
    }
    status, err := strconv.ParseBool(params[0])
    if err != nil {
        return nil, err
    }
    thoseNodeKey, thoseHostAddr := params[1], params[2]
    return &dcutcp.Response1013{ID:1013, Status:status, ThoseNodeKey: thoseNodeKey, Address: thoseHostAddr}, nil
}

func SplitCommand777Params(params []string) (*dcutcp.Command777, error) {
    if len(params) != 1 {
        return nil, errors.New(fmt.Sprintf("bad 777 request [%#v]", params))
    }
    status, err := strconv.ParseBool(params[0])
    if err != nil {
        return nil, err
    }
    return &dcutcp.Command777{ID:777, Status:status}, nil
}

func SplitCommand88Params(params []string) (*dcutcp.Command88, error) {
    if len(params) != 2 {
        return nil, errors.New(fmt.Sprintf("bad 88 request [%#v]", params))
    }
    return &dcutcp.Command88{ID:88, ThoseNodeKey: params[0], HostAddr:params[1]}, nil
}
func SplitTCPParams(params string) []string {
    out := make([]string, 0, 4)
    scanner := bufio.NewScanner(strings.NewReader(params))
    scanner.Split(SplitterFunc('\t'))
    for scanner.Scan() {
        out = append(out, scanner.Text())
    }
    return out
}

func ScanString(data string, delimiter byte) []string {
    out := make([]string, 0)
    scanner := bufio.NewScanner(strings.NewReader(data))
    scanner.Split(SplitterFunc(delimiter))
    for scanner.Scan() {
        out = append(out, scanner.Text())
    }
    return out
}

func RemovePortFromAddressString(address string) string {
    return strings.TrimSuffix(strings.TrimRightFunc(address,
        func(r rune) bool {
            return unicode.IsDigit(r)
        }), ":")
}