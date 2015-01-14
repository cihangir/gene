// Package constants generates the constant variables for a model/schema
package constants

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/writers"
)

// Generate generates the constants for given schema/model
func Generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("constants.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(ConstantsTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := temp.ExecuteTemplate(&buf, "constants.tmpl", s); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}

// ConstantsTemplate holds the template for the constant generation
var ConstantsTemplate = `
{{range $key, $value := .Properties}}
    {{if len $value.Enum}}
        // {{DepunctWithInitialUpper $key}} holds the predefined enums
        var {{DepunctWithInitialUpper $key}}  = struct {
        {{range $defKey, $val := $value.Enum}}
            {{DepunctWithInitialUpper $val}} string
        {{end}}
        }{
        {{range $defKey, $val := $value.Enum}}
            {{DepunctWithInitialUpper $val}}: "{{$val}}",
        {{end}}
        }
    {{end}}
{{end}}
`
