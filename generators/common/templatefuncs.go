package common

import (
	"strings"
	"text/template"

	"github.com/cihangir/gene/stringext"
)

var TemplateFuncs = template.FuncMap{
	"Pointerize":              stringext.Pointerize,
	"ToLower":                 strings.ToLower,
	"ToLowerFirst":            stringext.ToLowerFirst,
	"ToUpperFirst":            stringext.ToUpperFirst,
	"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
	"DepunctWithInitialLower": stringext.DepunctWithInitialLower,
	"Equal":                   stringext.Equal,
}
