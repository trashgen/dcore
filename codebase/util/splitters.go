package util

import (
    "fmt"
    "log"
    "bufio"
    "bytes"
    "errors"
    "strconv"
    "strings"
    "unicode"
    dchttputil "dcore/codebase/util/http"
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

///////////////////////////// START TCP SPLITTERS /////////////////////////////

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

func SplitPacket1013RequestParams(params []string) (*dchttputil.Request1013, error) {
    if len(params) != 1 {
        return nil, errors.New(fmt.Sprintf("bad 1013 request [%#v]", params))
    }
    return &dchttputil.Request1013{ID:1013, Key:params[0]}, nil
}

// TODO : Переделать по аналогии с SplitPacket1013RequestParams
func SplitPacket1013Response(data string) (*dchttputil.Response1013, error) {
    params := SplitTCPParams(strings.TrimSuffix(data, "\n"))
    if len(params) != 4 {
        return nil, errors.New(fmt.Sprintf("bad 1013 request [%s]", data))
    }

    id, err := strconv.Atoi(params[0])
    if err != nil {
        return nil, err
    }

    status, err := strconv.ParseBool(params[1])
    if err != nil {
        return nil, err
    }

    key, address := params[2], params[3]

    log.Printf("SplitPacket1013Response [%d] [%s] [%s]\n", id, key, address)
    return &dchttputil.Response1013{ID:id, Status:status, Key:key, Address:address}, nil
}

func SplitPacket777RequestParams(params []string) (*dchttputil.Request777, error) {
    if len(params) != 1 {
        return nil, errors.New(fmt.Sprintf("bad 777 request [%#v]", params))
    }
    status, err := strconv.ParseBool(params[0])
    if err != nil {
        return nil, err
    }
    return &dchttputil.Request777{ID:777, Status:status}, nil
}

func SplitPacket88RequestParams(params []string) (*dchttputil.Request88, error) {
    if len(params) != 1 {
        return nil, errors.New(fmt.Sprintf("bad 777 request [%#v]", params))
    }
    return &dchttputil.Request88{ID:777, Addr:params[0]}, nil
}

////////////////////////////// END TCP SPLITTERS ///////////////////////////////

///////////////////////////// START HTTP SPLITTERS ////////////////////////////

func SplitRequestReg(paramName string, queryParams string) (*dchttputil.RequestReg, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New(fmt.Sprintf("bad 'Reg' port param key [%s]", paramName))
    }

    port, err := strconv.Atoi(value)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("bad 'Reg' port value (not int) [%s]", value))
    }

    return &dchttputil.RequestReg{Port:port}, nil
}

func SplitRequestLook(paramName string, queryParams string) (*dchttputil.RequestLook, error) {
    if len(queryParams) == 0 {
        return &dchttputil.RequestLook{Count:0}, nil
    }

    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Look' Count param key")
    }

    count, err := strconv.Atoi(value)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("bad 'Look' 'count' param value [%s]", value))
    }

    return &dchttputil.RequestLook{Count:count}, nil
}

func SplitRequestCheck(paramName string, queryParams string) (*dchttputil.RequestCheck, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Check' 'key' param key")
    }

    return &dchttputil.RequestCheck{Key:value}, nil
}

func SplitRequestPoints(paramName string, queryParams string) (*dchttputil.RequestPoints, error) {
    if len(queryParams) == 0 {
        return &dchttputil.RequestPoints{Count:0}, nil
    }
    
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Points' 'count' param key")
    }

    count, err := strconv.Atoi(value)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("bad 'Points' 'count' param value [%s]", value))
    }
    
    return &dchttputil.RequestPoints{Count:count}, nil
}

func SplitRequestRemove(paramName string, queryParams string) (*dchttputil.RequestRemove, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Remove' 'key' param key")
    }
    
    return &dchttputil.RequestRemove{Key:value}, nil
}

/////////////////////////////// END HTTP SPLITTERS /////////////////////////////

func SplitQueryParams(queryParams string) map[string]string {
    out := make(map[string]string)
    scanner := bufio.NewScanner(strings.NewReader(queryParams))
    scanner.Split(SplitterFunc('&'))
    for scanner.Scan() {
        pair := strings.Split(scanner.Text(), "=")
        if len(pair) != 2 {
            log.Fatalf("Problem with query params: [%s]\n", scanner.Text())
        }

        out[pair[0]] = pair[1]
    }

    return out
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


func ScanString(data string, delimeter byte) []string {
    out := make([]string, 0)
    scanner := bufio.NewScanner(strings.NewReader(data))
    scanner.Split(SplitterFunc(delimeter))
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