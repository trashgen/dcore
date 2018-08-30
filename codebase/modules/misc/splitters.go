package misc

import (
    "fmt"
    "log"
    "bufio"
    "bytes"
    "errors"
    "strconv"
    "strings"
    dcutil "dcore/codebase/util"
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

func SplitRequestReg(paramName string, queryParams string) (*dcutil.RequestReg, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New(fmt.Sprintf("bad 'Reg' Address param key [%s]", paramName))
    }

    return &dcutil.RequestReg{Address:value}, nil
}

func SplitRequestLook(paramName string, queryParams string) (*dcutil.RequestLook, error) {
    if len(queryParams) == 0 {
        return &dcutil.RequestLook{Count:0}, nil
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

    return &dcutil.RequestLook{Count:count}, nil
}

func SplitRequestCheck(paramName string, queryParams string) (*dcutil.RequestCheck, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Check' 'key' param key")
    }
    
    return &dcutil.RequestCheck{Key:value}, nil
}

func SplitRequestPoints(paramName string, queryParams string) (*dcutil.RequestPoints, error) {
    if len(queryParams) == 0 {
        return &dcutil.RequestPoints{Count:0}, nil
    }
    
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Points' Count param key")
    }
    
    count, err := strconv.Atoi(value)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("bad 'Points' 'count' param value [%s]", value))
    }
    
    return &dcutil.RequestPoints{Count:count}, nil
}

func SplitRequestRemove(paramName string, queryParams string) (*dcutil.RequestRemove, error) {
    paramsMap := SplitQueryParams(queryParams)
    value, ok := paramsMap[paramName]
    if ! ok {
        return nil, errors.New("bad 'Remove' 'key' param key")
    }
    
    return &dcutil.RequestRemove{Key:value}, nil
}

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