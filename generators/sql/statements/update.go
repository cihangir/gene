package statements

import (
	"bytes"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// GenerateUpdate generates the update sql statement for the given schema
func GenerateUpdate(context *common.Context, s *schema.Schema) ([]byte, error) {
	temp := template.New("update_statement.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(UpdateStatementTemplate); err != nil {
		return nil, err
	}

	data := struct {
		Schema *schema.Schema
	}{
		Schema: s,
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "update_statement.tmpl", data); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// UpdateStatementTemplate holds the template for the update sql statement generator
var UpdateStatementTemplate = `
{{$title := Pointerize .Schema.Title}}
// GenerateUpdateSQL generates plain update sql statement for the given {{DepunctWithInitialUpper .Schema.Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Schema.Title}}) GenerateUpdateSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Update({{$title}}.TableName())

    {{range $key, $value := SortedSchema .Schema.Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $value.Title}}.IsZero(){
                psql = psql.Set("{{ToFieldName $value.Title}}", {{$title}}.{{DepunctWithInitialUpper $value.Title}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $value.Title}} != "" {
                psql = psql.Set("{{ToFieldName $value.Title}}", {{$title}}.{{DepunctWithInitialUpper $value.Title}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $value.Title}} != false {
                psql = psql.Set("{{ToFieldName $value.Title}}", {{$title}}.{{DepunctWithInitialUpper $value.Title}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $value.Title}}) != float64(0) {
                psql = psql.Set("{{ToFieldName $value.Title}}", {{$title}}.{{DepunctWithInitialUpper $value.Title}})
            }
        {{end}}
    {{end}}

    {{/* TODO get this ID section from the primary key*/}}
    return psql.Where("{{ToFieldName "Id"}} = ?", {{$title}}.{{DepunctWithInitialUpper "ID"}}).ToSql()
}
`
