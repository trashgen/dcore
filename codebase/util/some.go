package util

import (
    "bufio"
    "strings"
    dcmisc "dcore/codebase/modules/misc"
    "unicode"
)

func ScanString(data string, delimeter byte) []string {
    out := make([]string, 0)
    scanner := bufio.NewScanner(strings.NewReader(data))
    scanner.Split(dcmisc.SplitterFunc(delimeter))
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