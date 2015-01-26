package statements

import (
	"bytes"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// GenerateCreate generates the create sql statement for the given schema
func GenerateCreate(s *schema.Schema) ([]byte, error) {
	temp := template.New("create_statement.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(CreateStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "create_statement.tmpl", s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// CreateStatementTemplate holds the template for the create sql statement generator
var CreateStatementTemplate = `
// GenerateCreateSQL generates plain sql for the given {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) GenerateCreateSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Insert({{$title}}.TableName())
    columns := make([]string, 0)
    values := make([]interface{}, 0)

    {{range $key, $value := .Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $key}}.IsZero(){
                columns = append(columns, "{{ToFieldName $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != "" {
                columns = append(columns, "{{ToFieldName $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != false {
                columns = append(columns, "{{ToFieldName $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $key}}) != float64(0) {
                columns = append(columns, "{{ToFieldName $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{end}}
    {{end}}
    return psql.Columns(columns...).Values(values...).ToSql()
}
`
