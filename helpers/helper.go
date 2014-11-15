package helpers

import "regexp"

var NewLinesRegex = regexp.MustCompile(`(?m:\s*$)`)
