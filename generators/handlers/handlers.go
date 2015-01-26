// Package handlers creates the handlers for the gene
package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

// Generate commands the handler generation
func Generate(rootPath string, s *schema.Schema) error {
	for _, def := range s.Definitions {
		err := GenerateAPI(rootPath, s.Title, def.Title)
		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateAPI generates and writes the api files
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
