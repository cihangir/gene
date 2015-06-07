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

type generator struct{}

func New() *generator {
	return &generator{}
}

func (g *generator) Name() string {
	return "errors"
}

var PathForErrors = "%sworkers/%s/errors/%s.go"

// Generate generates and writes the errors of the schema
func (g *generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(s.Title)
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

		f, err := generate(def)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf(
			PathForErrors,
			context.Config.Target,
			moduleName,
			strings.ToLower(s.Title),
		)

		outputs = append(outputs, common.Output{
			Content: f,
			Path:    path,
		})

	}

	return outputs, nil
}

func generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("errors.tmpl").Funcs(common.TemplateFuncs)
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
