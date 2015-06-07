// Package errors generates the common errors for the modules
package errors

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Name() string {
	return "errors"
}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	outputs := make([]common.Output, 0)

	for _, def := range s.Definitions {
		// create models only for objects
		if def.Type != nil {
			if t, ok := def.Type.(string); ok {
				if t != "object" {
					continue
				}
			}
		}

		f, err := generate(context, def)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf(
			"%s/%s.go",
			context.Config.Target,
			strings.ToLower(def.Title),
		)

		outputs = append(outputs, common.Output{
			Content: f,
			Path:    path,
		})

	}

	return outputs, nil
}

func generate(context *common.Context, s *schema.Schema) ([]byte, error) {
	temp := template.New("errors.tmpl").Funcs(context.TemplateFuncs)
	_, err := temp.Parse(ErrorsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "errors.tmpl", s); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}

// ErrorsTemplate holds the template for the errors package
var ErrorsTemplate = `
package errs
var (
{{$moduleName := ToUpperFirst .Title}}
{{range $key, $value := .Properties}}
    Err{{$moduleName}}{{ToUpperFirst $key}}NotSet = errors.New("{{$moduleName}}.{{ToUpperFirst $key}} not set")
{{end}}
)
`
