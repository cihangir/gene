package rows

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Name() string {
	return "scanner-row"
}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	outputs := make([]common.Output, 0)
	for _, key := range schema.SortedKeys(s.Definitions) {
		def := s.Definitions[key]
		output, err := GenerateScanner(context.Config.Target, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

// GenerateScanner generates and writes the api files
func GenerateScanner(rootPath string, s *schema.Schema) (common.Output, error) {
	api, err := generate(s)
	if err != nil {
		return common.Output{}, err
	}

	path := fmt.Sprintf("%s%s_rowscanner.go", rootPath, strings.ToLower(s.Title))

	return common.Output{
		Content: api,
		Path:    path,
	}, nil
}

// RowScannerTemplate provides the template for rowscanner of models
var RowScannerTemplate = `
{{$schema := .Schema}}
{{$title := $schema.Title}}
package models

func ({{Pointerize $title}} *{{$title}}) RowsScan(rows *sql.Rows, dest interface{}) error {
    if rows == nil {
        return nil
    }

    var records []*{{ToUpperFirst $title}}
    for rows.Next() {
        m := New{{ToUpperFirst $title}}()
        err := rows.Scan(
        {{range $n, $p := $schema.Properties}} &m.{{DepunctWithInitialUpper $p.Title}},
        {{end}} )
        if err != nil {
            return err
        }
        records = append(records, m)
    }

    if err := rows.Err(); err != nil {
        return err
    }

    *(dest.(*[]*{{ToUpperFirst $title}})) = records

    return nil
}
`

// Generate generates the rowscanner for given schema/model
func generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("rowscanner.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(RowScannerTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema *schema.Schema
	}{
		Schema: s,
	}

	if err := temp.ExecuteTemplate(&buf, "rowscanner.tmpl", data); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
