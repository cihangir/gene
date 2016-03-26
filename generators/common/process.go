package common

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
)

type Op struct {
	Name         string
	Template     string
	PathFunc     func(context *Context, def *schema.Schema, moduleName string) string
	Clear        bool
	DoNotFormat  bool
	FormatSource bool
	// TemplateFuncs template.FuncMap
}

func Proces(o *Op, req *Req, res *Res) error {
	if req == nil || req.Context == nil || req.Context.Config == nil {
		return nil
	}

	if !IsIn(o.Name, req.Context.Config.Generators...) {
		return nil
	}

	tmpl := template.New("dockerfile.tmpl").Funcs(req.Context.TemplateFuncs)
	if _, err := tmpl.Parse(o.Template); err != nil {
		return err
	}
	// s := req.Schema
	// c := req.Context.Config

	moduleName := req.Context.ModuleNameFunc(req.Schema.Title)

	outputs := make([]Output, 0)

	for _, def := range SortedObjectSchemas(req.Schema.Definitions) {
		data := struct {
			ModuleName string
			Schema     *schema.Schema
		}{
			ModuleName: moduleName,
			Schema:     def,
		}

		var buf bytes.Buffer

		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}

		var content []byte
		var err error
		if o.Clear {
			content, err = writers.Clear(buf)
			if err != nil {
				return err
			}
		} else {
			content = buf.Bytes()
		}

		if o.FormatSource {
			content, err = format.Source(content)
			if err != nil {
				return err
			}
		}

		outputs = append(outputs, Output{
			Content:     content,
			Path:        o.PathFunc(req.Context, def, moduleName),
			DoNotFormat: o.DoNotFormat,
		})
	}

	res.Output = outputs
	return nil
}
