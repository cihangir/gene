package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/format"
	"strings"
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

	if req.Schema == nil {
		if req.SchemaStr == "" {
			return errors.New("both schema and string schema is not set")
		}

		s := &schema.Schema{}
		if err := json.Unmarshal([]byte(req.SchemaStr), s); err != nil {
			return err
		}

		req.Schema = s.Resolve(nil)
	}

	settings, ok := req.Schema.Generators.Get(o.Name)
	if !ok {
		settings = schema.Generator{}
	}

	tmpl := template.New("dockerfile.tmpl").Funcs(TemplateFuncs)
	if _, err := tmpl.Parse(o.Template); err != nil {
		return err
	}

	moduleName := strings.ToLower(req.Schema.Title)

	outputs := make([]Output, 0)

	for _, def := range SortedObjectSchemas(req.Schema.Definitions) {
		data := struct {
			ModuleName string
			Schema     *schema.Schema
			Settings   *schema.Generator
		}{
			ModuleName: moduleName,
			Schema:     def,
			Settings:   &settings,
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
