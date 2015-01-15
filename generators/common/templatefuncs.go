// Package common provides common operation helpers to the generators
package common

import (
	"strings"
	"text/template"

	"github.com/cihangir/gene/stringext"
)

// TemplateFuncs provides utility functions for template operations
var TemplateFuncs = template.FuncMap{
	"Pointerize":              stringext.Pointerize,
	"ToLower":                 strings.ToLower,
	"ToLowerFirst":            stringext.ToLowerFirst,
	"ToUpperFirst":            stringext.ToUpperFirst,
	"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
	"DepunctWithInitialLower": stringext.DepunctWithInitialLower,
	"Equal":                   stringext.Equal,
	"ToFieldName":             stringext.ToFieldName,
}
