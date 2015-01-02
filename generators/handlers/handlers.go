package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, name string) error {
	err := GenerateAPI(rootPath, name)
	if err != nil {
		return err
	}

	return GenerateHandlers(rootPath, name)

}

func GenerateAPI(rootPath string, name string) error {
	temp := template.New("api.tmpl")
	temp.Funcs(template.FuncMap{
		"ToLowerFirst": stringext.ToLowerFirst,
	})

	_, err := temp.Parse(APITemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "api.tmpl", name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sgene/modules/%s/api/%s.go",
		rootPath,
		stringext.ToLowerFirst(name),
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
		"%sgene/modules/%s/handlers/%s.go",
		rootPath,
		stringext.ToLowerFirst(name),
		stringext.ToLowerFirst(name),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
