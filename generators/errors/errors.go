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

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(req *common.Req, res *common.Res) error {
	context := req.Context
	s := req.Schema

	if context == nil || context.Config == nil {
		return nil
	}

	if !common.IsIn("errors", context.Config.Generators...) {
		return nil
	}

	temp := template.New("errors.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(ErrorsTemplate); err != nil {
		return err
	}

	outputs := make([]common.Output, 0)

	for _, def := range common.SortedObjectSchemas(s.Definitions) {
		data := struct {
			Schema *schema.Schema
		}{
			Schema: def,
		}

		var buf bytes.Buffer

		if err := temp.ExecuteTemplate(&buf, "errors.tmpl", data); err != nil {
			return err
		}

		f, err := writers.Clear(buf)
		if err != nil {
			return err
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

	res.Output = outputs
	return nil
}
