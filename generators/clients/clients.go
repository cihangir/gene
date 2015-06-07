// Package clients generates clients for the generated api
package clients

import (
	"bytes"
	"fmt"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

func New() *Generator {
	return &Generator{}
}

func (c *Generator) Name() string {
	return "clients"
}

// Generate generates the client package for given schema
func (c *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(s.Title)
	outputs := make([]common.Output, 0)

	for _, key := range schema.SortedKeys(s.Definitions) {
		def := s.Definitions[key]

		if def.Type != nil {
			if t, ok := def.Type.(string); ok {
				if t != "object" {
					continue
				}
			}
		}

		f, err := c.generate(context, moduleName, def)
		if err != nil {
			return outputs, err
		}

		path := fmt.Sprintf(
			"%s%s/clients/%s.go",
			context.Config.Target,
			moduleName,
			context.FileNameFunc(def.Title),
		)

		outputs = append(outputs, common.Output{Content: f, Path: path})
	}

	return outputs, nil
}

func (c *Generator) generate(context *common.Context, moduleName string, s *schema.Schema) ([]byte, error) {
	tmpl := template.New("clients.tmpl").Funcs(context.TemplateFuncs)
	if _, err := tmpl.Parse(ClientsTemplate); err != nil {
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

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
