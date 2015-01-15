package statements

import (
	"bytes"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
)

// GenerateTableName generates a simple table name getter function
func GenerateTableName(s *schema.Schema) ([]byte, error) {
	temp := template.New("table_name_statement.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(TableNameTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "table_name_statement.tmpl", s); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

// TableNameTemplate holds the template for the TableName function
var TableNameTemplate = `
// TableName returns the table name for {{DepunctWithInitialUpper .Title}}
{{$title := Pointerize .Title}}
func ({{$title}} *{{DepunctWithInitialUpper .Title}}) TableName() string {
    return "{{DepunctWithInitialLower .Title}}"
}
`
