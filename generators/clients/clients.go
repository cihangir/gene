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

func NewClient(context *config.Context, schema *schema.Schema) (*client, error) {
	// create a template to process the clients
	tmpl := template.New("clients.tmpl").Funcs(context.TemplateFuncs)
	if _, err := tmpl.Parse(ClientsTemplate); err != nil {
		return nil, err
	}

	c := &client{
		context:  context,
		schema:   schema,
		template: tmpl,
	}

	return c, nil
}

// Generate generates the client package for given schema
func (c *client) Generate() ([]common.Output, error) {
	moduleName := c.context.ModuleNameFunc(c.schema.Title)

	keys := schema.SortedKeys(c.schema.Definitions)
	outputs := make([]common.Output, len(keys))

	for i, key := range keys {
		def := c.schema.Definitions[key]

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
			c.context.Config.Target,
			moduleName,
			c.context.FileNameFunc(def.Title),
		)

		outputs[i] = common.Output{
			Content: f,
			Path:    path,
		}
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
