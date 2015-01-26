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

// Generate generates and writes the errors of the schema
func Generate(rootPath string, s *schema.Schema) error {
	data, err := generate(s)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/errors/%s.go",
		rootPath,
		strings.ToLower(s.Title),
		strings.ToLower(s.Title),
	)

	return writers.WriteFormattedFile(path, data)
}

func generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("errors.tmpl").Funcs(common.TemplateFuncs)
	_, err := temp.Parse(ErrorsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "errors.tmpl", s); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}
