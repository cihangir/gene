package statements

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
)

func GenerateCreate(s *schema.Schema) ([]byte, error) {
	temp := template.New("create_statement.tmpl")

	temp.Funcs(template.FuncMap{
		"Pointerize":              stringext.Pointerize,
		"DepunctWithInitialUpper": stringext.DepunctWithInitialUpper,
		"Equal":                   stringext.Equal,
		"ToSnake":                 stringext.ToSnake,
	})
	_, err := temp.Parse(CreateStatementTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "create_statement.tmpl", s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var CreateStatementTemplate = `
// GenerateCreateSQL generates plain sql for the given {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) GenerateCreateSQL() (string, []interface{}, error) {
    psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    psql = psql.Insert({{$title}}.TableName())
    columns := make([]string, 0)
    values := make([]interface{}, 0)

    {{range $key, $value := .Properties}}
        {{/* handle strings */}}
        {{if Equal "string" $value.Type}}
            {{/* strings can have special formatting */}}
            {{if Equal "date-time" $value.Format}}
            if !{{$title}}.{{DepunctWithInitialUpper $key}}.IsZero(){
                columns = append(columns, "{{ToSnake $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{else}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != "" {
                columns = append(columns, "{{ToSnake $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
            {{end}}

        {{else if Equal "boolean" $value.Type}}
            if {{$title}}.{{DepunctWithInitialUpper $key}} != false {
                columns = append(columns, "{{ToSnake $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{else if Equal "number" $value.Type}}
            if float64({{$title}}.{{DepunctWithInitialUpper $key}}) != float64(0) {
                columns = append(columns, "{{ToSnake $key}}")
                values = append(values, {{$title}}.{{DepunctWithInitialUpper $key}})
            }
        {{end}}
    {{end}}
    return psql.Columns(columns...).Values(values...).ToSql()
}
`
