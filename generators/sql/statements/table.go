package statements

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
)

func GenerateTableName(s *schema.Schema) ([]byte, error) {
	temp := template.New("table_name_statement.tmpl").Funcs(common.TemplateFuncs)
	_, err := temp.Parse(TableNameTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "table_name_statement.tmpl", s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var TableNameTemplate = `
// TableName returns the table name for a given struct
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) TableName() string {
    return "{{DepunctWithInitialLower .Title}}"
}
`
