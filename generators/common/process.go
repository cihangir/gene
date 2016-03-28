package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/format"
	"strings"
	"text/template"

	"github.com/cihangir/gene/utils"
	"github.com/cihangir/schema"
)

type PostProcessor func([]byte) []byte

type Op struct {
	Name           string
	Template       string
	PathFunc       func(data *TemplateData) string
	Clear          bool
	DoNotFormat    bool
	FormatSource   bool
	PostProcessors []PostProcessor
	// TemplateFuncs template.FuncMap
}

type TemplateData struct {
	ModuleName string
	Schema     *schema.Schema
	Settings   *schema.Generator
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

	settings.SetNX("rootPathPrefix", o.Name)
	rootPathPrefix := settings.Get("rootPathPrefix").(string)
	fullPathPrefix := req.Context.Config.Target + rootPathPrefix + "/"
	settings.Set("fullPathPrefix", fullPathPrefix)

	tmpl := template.New("template").Funcs(TemplateFuncs)
	if _, err := tmpl.Parse(o.Template); err != nil {
		return err
	}

	moduleName := strings.ToLower(req.Schema.Title)

	outputs := make([]Output, 0)

	for _, def := range SortedObjectSchemas(req.Schema.Definitions) {
		data := &TemplateData{
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
			content, err = utils.Clear(buf)
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
			Path:        o.PathFunc(data),
			DoNotFormat: o.DoNotFormat,
		})
	}
	res.Output = outputs
	return nil
}

func ProcessSingle(o *Op, def *schema.Schema, settings schema.Generator) ([]byte, error) {
	temp := template.New("single").Funcs(TemplateFuncs)
	if _, err := temp.Parse(o.Template); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema     *schema.Schema
		Settings   schema.Generator
		Properties []*schema.Schema
	}{
		Schema:     def,
		Settings:   settings,
		Properties: schema.SortedSchema(def.Properties),
	}

	if err := temp.ExecuteTemplate(&buf, "single", data); err != nil {
		return nil, err
	}

	b := buf.Bytes()

	for _, processor := range o.PostProcessors {
		b = processor(b)
	}

	return b, nil
}
