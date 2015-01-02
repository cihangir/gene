package errors

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/generators/validators"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {
	temp := template.New("errors.tmpl")

	temp.Funcs(template.FuncMap{
		"GenerateValidator": validators.GenerateValidator,
		"Pointerize":        stringext.Pointerize,
	})

	_, err := temp.Parse(ErrorsTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "errors.tmpl", s); err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/errors/%s.go",
		rootPath,
		stringext.ToLowerFirst(s.Title),
		stringext.ToLowerFirst(s.Title),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
