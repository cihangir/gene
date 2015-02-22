// Package constructors generates the constructors for given schema/model
package constructors

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
)

// Generate generates the constructors for given schema/model
func Generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("constructors.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(ConstructorsTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "constructors.tmpl", s); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}

// ConstructorsTemplate provides the template for constructors of models
var ConstructorsTemplate = `
// New{{DepunctWithInitialUpper .Title}} creates a new {{DepunctWithInitialUpper .Title}} struct with default values
func New{{DepunctWithInitialUpper .Title}}() *{{DepunctWithInitialUpper .Title}} {
    return &{{DepunctWithInitialUpper .Title}}{
        {{range $key, $value := .Properties}}
            {{/* only process if default value is set */}}
            {{if $value.Default}}
                {{/* handle strings */}}
                {{if Equal "string" $value.Type}}
                    {{/* if property is enum, handle them accordingly */}}
                    {{if len $value.Enum}}
                        {{DepunctWithInitialUpper $key}}: {{DepunctWithInitialUpper $key}}.{{DepunctWithInitialUpper $value.Default}},
                    {{else}}
                        {{/* strings can have special formatting */}}
                        {{/* no-matter what value set for a date-time field, set UTC Now */}}
                        {{if Equal "date-time" $value.Format}}
                            {{DepunctWithInitialUpper $key}}: time.Now().UTC(),
                        {{else}}
                            {{DepunctWithInitialUpper $key}}: "{{$value.Default}}",
                        {{end}}
                    {{end}}
                {{else}}
                    {{/* for boolean, numbers.. */}}
                    {{DepunctWithInitialUpper $key}}: {{$value.Default}},
                {{end}}
            {{end}}
        {{end}}
    }
}
`
