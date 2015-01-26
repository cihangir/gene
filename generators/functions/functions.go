package functions

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
)

// Generate generates and writes the errors of the schema
func Generate(rootPath string, s *schema.Schema) error {
	a, err := generate(s)
	fmt.Println("string(a), err-->", string(a), err)

	return nil
}

// ConstructorsTemplate provides the template for constructors of models
var ConstructorsTemplate = `
{{$title := .Title}}
{{range $defKey, $defValue := .Definitions}}
    {{range $funcKey, $funcValue := $defValue.Functions}}
    func ({{Pointerize $title}} *{{$title}}) {{$funcKey}}(req *{{Argumentize $funcValue.Properties.incoming}}, res *{{Argumentize $funcValue.Properties.outgoing}}) error {

    }
    {{end}}
{{end}}
`

// Generate generates the constructors for given schema/model
func generate(s *schema.Schema) ([]byte, error) {
	temp := template.New("constructors.tmpl").Funcs(common.TemplateFuncs)

	if _, err := temp.Parse(ConstructorsTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := temp.ExecuteTemplate(&buf, "constructors.tmpl", s); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}
