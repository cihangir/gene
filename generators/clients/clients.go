// Package clients generates clients for the generated api
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

// PathForClient holds the to be formatted string for the path of the client
var PathForClient = "%sworkers/%s/clients/%s.go"

// Generate generates the client package for given schema
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

		f, err := generate(moduleName, def)
		if err != nil {
			return err
		}

		path := fmt.Sprintf(PathForClient, rootPath, moduleName,
			stringext.ToLower(def.Title),
		)

		if err := writers.WriteFormattedFile(path, f); err != nil {
			return err
		}

	}

	return nil
}

func generate(moduleName string, s *schema.Schema) ([]byte, error) {
	// create a template to process the clients
	temp := template.New("clients.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(ClientsTemplate); err != nil {
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

	if err := temp.ExecuteTemplate(&buf, "clients.tmpl", data); err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
