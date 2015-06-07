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

type generator struct{}

func New() *generator {
	return &generator{}
}

var PathForFunctions = "%sworkers/%s/api/%s.go"

func (g *generator) Name() string {
	return "functions"
}

// Generate generates and writes the errors of the schema
func (g *generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(s.Title)
	outputs := make([]common.Output, 0)

	keys := schema.SortedKeys(s.Definitions)
	for _, key := range keys {
		def := s.Definitions[key]
		output, err := GenerateAPI(
			context.Config.Target,
			moduleName,
			def,
		)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output)
	}

	return outputs, nil
}

// GenerateAPI generates and writes the api files
func GenerateAPI(rootPath string, moduleName string, s *schema.Schema) (common.Output, error) {
	api, err := generate(moduleName, s)
	if err != nil {
		return common.Output{}, err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/api/%s.go",
		rootPath,
		moduleName,
		strings.ToLower(s.Title),
	)

	return common.Output{
		Content: api,
		Path:    path,
	}, nil
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
func generate(moduleName string, s *schema.Schema) ([]byte, error) {
	temp := template.New("constructors.tmpl").Funcs(common.TemplateFuncs)

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
