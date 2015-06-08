package functions

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
	return "functions"
}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(s.Title)
	outputs := make([]common.Output, 0)

	for _, def := range schema.SortedSchema(s.Definitions) {
		api, err := generate(context, moduleName, def)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf(
			"%s%s/api/%s.go",
			context.Config.Target,
			moduleName,
			strings.ToLower(s.Title),
		)

		outputs = append(outputs, common.Output{
			Content: api,
			Path:    path,
		})
	}

	return outputs, nil
}

// FunctionsTemplate provides the template for constructors of models
var FunctionsTemplate = `
{{$schema := .Schema}}
{{$title := $schema.Title}}

package {{ToLower .ModuleName}}api

// New creates a new local {{ToUpperFirst $title}} handler
func New{{ToUpperFirst $title}}() *{{ToUpperFirst $title}} { return &{{ToUpperFirst $title}}{} }

// {{ToUpperFirst $title}} is for holding the api functions
type {{ToUpperFirst $title}} struct{}

{{range $funcKey, $funcValue := $schema.Functions}}
func ({{Pointerize $title}} *{{$title}}) {{$funcKey}}(ctx context.Context, req *{{Argumentize $funcValue.Properties.incoming}}, res *{{Argumentize $funcValue.Properties.outgoing}}) error {
    return db.MustGetDB(ctx).{{$funcKey}}(models.New{{ToUpperFirst $title}}(), req, res)
}
{{end}}
`

// Generate generates the constructors for given schema/model
func generate(context *common.Context, moduleName string, s *schema.Schema) ([]byte, error) {
	temp := template.New("constructors.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(FunctionsTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		ModuleName string
		Schema     *schema.Schema
	}{
		ModuleName: moduleName,
		Schema:     s,
	}

	if err := temp.ExecuteTemplate(&buf, "constructors.tmpl", data); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
