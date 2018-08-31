package util

import (
    "bufio"
    "strings"
    dcmisc "dcore/codebase/modules/misc"
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