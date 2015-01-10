package constants

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
)

func Generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("constants.tmpl")

	temp.Funcs(template.FuncMap{
		"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
	})
	_, err := temp.Parse(ConstantsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "constants.tmpl", s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var ConstantsTemplate = `
{{range $key, $value := .Properties}}
    {{if len $value.Enum}}
        // {{DepunctWithInitialUpper $key}} holds the predefined enums
        const {{DepunctWithInitialUpper $key}}  = struct {
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
