package constructors

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/writers"
)

func Generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("constructors.tmpl").Funcs(common.TemplateFuncs)
	_, err := temp.Parse(ConstructorsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "constructors.tmpl", s)
	if err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}

var ConstructorsTemplate = `
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
