package stringext

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	camelcase = regexp.MustCompile(`(?m)[-.$/:_{}\s]`)
	acronyms  = regexp.MustCompile(`(Url|Http|Id|Io|Uuid|Api|Uri|Ssl|Cname|Oauth|Otp)$`)
)

func ToLowerFirst(ident string) string {
	r, n := utf8.DecodeRuneInString(ident)
	return string(unicode.ToLower(r)) + ident[n:]
}

func Pointerize(ident string) string {
	r, _ := utf8.DecodeRuneInString(ident)
	return string(unicode.ToLower(r))
}

func ToUpperFirst(ident string) string {
	r, n := utf8.DecodeRuneInString(ident)
	return string(unicode.ToUpper(r)) + ident[n:]
}

func JSONTag(n string, required bool) string {
	tags := []string{ToLowerFirst(n)}
	if !required {
		tags = append(tags, "omitempty")
	}

	return fmt.Sprintf("`json:\"%s\"`", strings.Join(tags, ","))
}

// Equal check if given two strings are same, used in templates
func Equal(a, b string) bool {
	return a == b
}

func AsComment(c string) string {
	var buf bytes.Buffer
	const maxLen = 70
	removeNewlines := func(s string) string {
		return strings.Replace(s, "\n", "\n// ", -1)
	}
	for len(c) > 0 {
		line := c
		if len(line) < maxLen {
			fmt.Fprintf(&buf, "// %s\n", removeNewlines(line))
			break
		}
		line = line[:maxLen]
		si := strings.LastIndex(line, " ")
		if si != -1 {
			line = line[:si]
		}
		fmt.Fprintf(&buf, "// %s\n", removeNewlines(line))
		c = c[len(line):]
		if si != -1 {
			c = c[1:]
		}
	}

	return buf.String()
}

func Contains(n string, r []string) bool {
	for _, r := range r {
		if r == n {
			return true
		}
	}

	return false
}

func DepunctWithInitialUpper(ident string) string {
	return Depunct(ident, true)
}

func DepunctWithInitialLower(ident string) string {
	return Depunct(ident, false)
}

func Depunct(ident string, initialCap bool) string {
	matches := camelcase.Split(ident, -1)
	for i, m := range matches {
		if initialCap || i > 0 {
			m = ToUpperFirst(m)
		}
		matches[i] = acronyms.ReplaceAllStringFunc(m, func(c string) string {
			if len(c) > 4 {
				return strings.ToUpper(c[:2]) + c[2:]
			}
			return strings.ToUpper(c)
		})
	}
	return strings.Join(matches, "")
}
