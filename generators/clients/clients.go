// Package clients generates clients for the generated api
package clients

import (
	"bytes"
	"fmt"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type client struct {
	context  *config.Context
	schema   *schema.Schema
	template *template.Template
}

func New() *client {
	// create a template to process the clients
	return &client{}
}

func (c *client) Name() string {
	return "clients"
}

// Generate generates the client package for given schema
func (c *client) Generate(context *config.Context, s *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(s.Title)
	keys := schema.SortedKeys(s.Definitions)
	outputs := make([]common.Output, 0)
	tmpl := template.New("clients.tmpl").Funcs(context.TemplateFuncs)

	c.context = context
	c.schema = s
	c.template = tmpl

	if _, err := tmpl.Parse(ClientsTemplate); err != nil {
		return nil, err
	}

	for _, key := range keys {
		def := s.Definitions[key]

		if def.Type != nil {
			if t, ok := def.Type.(string); ok {
				if t != "object" {
					continue
				}
			}
		}

		f, err := c.generate(moduleName, def)
		if err != nil {
			return outputs, err
		}

		path := fmt.Sprintf(PathForClient,
			context.Config.Target,
			moduleName,
			context.FileNameFunc(def.Title),
		)

		outputs = append(outputs, common.Output{
			Content: f,
			Path:    path,
		})
	}

	return outputs, nil
}

// PathForClient holds the to be formatted string for the path of the client
var PathForClient = "%sworkers/%s/clients/%s.go"

func (c *client) generate(moduleName string, s *schema.Schema) ([]byte, error) {

	var buf bytes.Buffer

	data := struct {
		ModuleName string
		Schema     *schema.Schema
	}{
		ModuleName: moduleName,
		Schema:     s,
	}

	if err := c.template.Execute(&buf, data); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
