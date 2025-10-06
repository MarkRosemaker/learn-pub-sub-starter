package util

import (
	"strings"
)

const Asterisk = "*"

func DotJoin(parts ...string) string {
	return strings.Join(parts, ".")
}
