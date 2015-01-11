package common

import (
	"strings"
	"text/template"

	"github.com/cihangir/gene/stringext"
)

var TemplateFuncs = template.FuncMap{
	"Pointerize":              stringext.Pointerize,
	"ToLowerFirst":            stringext.ToLowerFirst,
	"ToLower":                 strings.ToLower,
	"ToUpperFirst":            stringext.ToUpperFirst,
	"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
}
