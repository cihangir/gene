package clients

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
	moduleName := stringext.ToLowerFirst(s.Title)

	path := fmt.Sprintf(
		"%sworkers/%s/clients/%s.go",
		rootPath,
		moduleName,
		moduleName,
	)

	f, err := generateClient(s)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(path, f)
}

// Generate functions according to the schema.
func generateClient(s *schema.Schema) ([]byte, error) {
	temp := template.New("clients.tmpl")
	temp.Funcs(template.FuncMap{
		"GenerateValidator": validators.GenerateValidator,
		"Pointerize":        stringext.Pointerize,
		"ToLowerFirst":      stringext.ToLowerFirst,
	})

	_, err := temp.Parse(ClientsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "clients.tmpl", s.Title)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
