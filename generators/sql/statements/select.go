package statements

import (
	"bytes"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// GenerateSelect generates the select sql statement for the given schema
func GenerateSelect(s *schema.Schema) ([]byte, error) {
	temp := template.New("select_statement.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(SelectStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "select_statement.tmpl", s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// SelectStatementTemplate holds the template for the select sql statement generator
var SelectStatementTemplate = `
// GenerateSelectSQL generates plain select sql statement for the given {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) GenerateSelectSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From({{$title}}.TableName())

    columns := make([]string, 0)
    values := make([]interface{}, 0)

    {{range $key, $value := .Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $key}}.IsZero(){
                columns = append(columns, "{{ToFieldName $key}} = ?")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != "" {
                columns = append(columns, "{{ToFieldName $key}} = ?")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != false {
                columns = append(columns, "{{ToFieldName $key}} = ?")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $key}}) != float64(0) {
                columns = append(columns, "{{ToFieldName $key}} = ?")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{end}}
    {{end}}

    return psql.Where(strings.Join(columns, " AND "), values...).ToSql()
}
`
