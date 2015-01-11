package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"go/format"

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
}

func GenerateAPI(rootPath string, moduleName string, name string) error {
	api, err := generate(moduleName, name)
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

	return writers.WriteFormattedFile(path, api)
}

func generate(moduleName string, name string) ([]byte, error) {
	temp := template.New("api.tmpl")
	temp.Funcs(common.TemplateFuncs)

	_, err := temp.Parse(APITemplate)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return format.Source(buf.Bytes())
}
