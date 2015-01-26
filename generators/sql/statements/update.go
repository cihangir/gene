package statements

import (
	"bytes"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// GenerateUpdate generates the update sql statement for the given schema
func GenerateUpdate(s *schema.Schema) ([]byte, error) {
	temp := template.New("update_statement.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(UpdateStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "update_statement.tmpl", s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// UpdateStatementTemplate holds the template for the update sql statement generator
var UpdateStatementTemplate = `
// GenerateUpdateSQL generates plain update sql statement for the given {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) GenerateUpdateSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update({{$title}}.TableName())

    {{range $key, $value := .Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $key}}.IsZero(){
                psql = psql.Set("{{ToFieldName $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != "" {
                psql = psql.Set("{{ToFieldName $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != false {
                psql = psql.Set("{{ToFieldName $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $key}}) != float64(0) {
                psql = psql.Set("{{ToFieldName $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{end}}
    {{end}}

    {{/* TODO get this ID section from the primary key*/}}
    return psql.Where("{{ToFieldName "Id"}} = ?", {{$title}}.{{DepunctWithInitialUpper "ID"}}).ToSql()
}
`
