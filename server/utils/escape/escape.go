// Package escape string is used for prevent sql injections
package escape

import (
	"strings"
)

func MakeSafe(value string) string {
    replace := map[string]string{"\\":"\\\\", "'":`\'`, "\\0":"\\\\0", "\n":"\\n", "\r":"\\r", `"`:`\"`, "\x1a":"\\Z"}

    for b, a := range replace {
        value = strings.Replace(value, b, a, -1)
    }

    return value
}
