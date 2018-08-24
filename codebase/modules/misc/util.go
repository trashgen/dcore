package misc

import "bytes"

func SplitTabs(data []byte, atEOF bool) (advance int, token []byte, err error) {
    if atEOF && len(data) == 0 {
        return 0, nil, nil
    }
    if i := bytes.IndexByte(data, '\t'); i >= 0 {
        return i + 1, data[0:i], nil
    }
    if atEOF {
        return len(data), data, nil
    }
    return 0, nil, nil
}