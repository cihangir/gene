package handlers

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {
	for _, def := range s.Definitions {
		err := GenerateAPI(rootPath, s.Title, def.Title)
		if err != nil {
			return err
		}
	}

	return nil
	// return GenerateHandlers(rootPath, name)
}

func GenerateAPI(rootPath string, moduleName string, name string) error {
	temp := template.New("api.tmpl")
	temp.Funcs(common.TemplateFuncs)

	_, err := temp.Parse(APITemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	data := struct {
		Name       string
		ModuleName string
	}{
		Name:       name,
		ModuleName: moduleName,
	}

	err = temp.ExecuteTemplate(&buf, "api.tmpl", data)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/%sapi/%s.go",
		rootPath,
		moduleName,
		moduleName,
		stringext.ToLowerFirst(name),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}

func GenerateHandlers(rootPath string, name string) error {
	temp := template.New("handlers.tmpl")
	temp.Funcs(template.FuncMap{
		"ToLowerFirst": stringext.ToLowerFirst,
	})

	_, err := temp.Parse(HandlersTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "handlers.tmpl", name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/handlers/%s.go",
		rootPath,
		stringext.ToLowerFirst(name),
		stringext.ToLowerFirst(name),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
