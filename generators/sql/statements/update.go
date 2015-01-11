package statements

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
)

func GenerateUpdate(s *schema.Schema) ([]byte, error) {
	temp := template.New("update_statement.tmpl")
	temp.Funcs(template.FuncMap{
		"Pointerize":              stringext.Pointerize,
		"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
		"Equal":                   stringext.Equal,
		"ToSnake":                 stringext.ToSnake,
		"DepunctWithInitialLower": stringext.DepunctWithInitialLower,
	})

	_, err := temp.Parse(UpdateStatementTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "update_statement.tmpl", s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var UpdateStatementTemplate = `
// GenerateUpdateSQL generates plain update sql statement for the given {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) GenerateUpdateSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    psql = psql.Update({{$title}}.TableName())

    {{range $key, $value := .Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $key}}.IsZero(){
                psql = psql.Set("{{ToSnake $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != "" {
                psql = psql.Set("{{ToSnake $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != false {
                psql = psql.Set("{{ToSnake $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $key}}) != float64(0) {
                psql = psql.Set("{{ToSnake $key}}", {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{end}}
    {{end}}

    {{/* TODO get this ID section from the primary key*/}}
    return psql.Where("{{ToSnake "Id"}} = ?", {{$title}}.{{DepunctWithInitialUpper "ID"}}).ToSql()
}
`
