package clients

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
	moduleName := stringext.ToLowerFirst(s.Title)

	for _, def := range s.Definitions {

		if def.Type != nil {
			if t, ok := def.Type.(string); ok {
				if t != "object" {
					continue
				}
			}
		}

		path := fmt.Sprintf(
			"%sworkers/%s/clients/%s.go",
			rootPath,
			moduleName,
			stringext.ToLowerFirst(def.Title),
		)

		f, err := generate(moduleName, def)
		if err != nil {
			return err
		}

		err = writers.WriteFormattedFile(path, f)
		if err != nil {
			return err
		}

	}

	return nil
}

// Generate functions according to the schema.
func generate(moduleName string, s *schema.Schema) ([]byte, error) {
	temp := template.New("clients.tmpl")
	temp.Funcs(common.TemplateFuncs)

	_, err := temp.Parse(ClientsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Name       string
		ModuleName string
	}{
		Name:       stringext.ToLowerFirst(s.Title),
		ModuleName: moduleName,
	}

	err = temp.ExecuteTemplate(&buf, "clients.tmpl", data)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
