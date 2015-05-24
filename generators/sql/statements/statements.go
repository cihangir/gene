package statements

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
)

type generator struct{}

func New() *generator {
	return &generator{}
}

var PathForStatements = "%smodels/%s_statements.go"

func (g *generator) Name() string {
	return "statements"
}

// Generate generates the basic CRUD statements for the models
func (g *generator) Generate(context *config.Context, schema *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(schema.Title)
	outputs := make([]common.Output, 0)

	for _, def := range schema.Definitions {
		// create models only for objects
		if def.Type != nil {
			if t, ok := def.Type.(string); ok {
				if t != "object" {
					continue
				}
			}
		}

		f, err := GenerateModelStatements(def)
		if err != nil {
			return outputs, err
		}

		path := fmt.Sprintf(
			PathForStatements,
			context.Config.Target,
			moduleName,
		)

		outputs = append(outputs, common.Output{
			Content: f,
			Path:    path,
		})

	}
	return outputs, nil
}

// GenerateModelStatements generates the CRUD statements for the model struct
func GenerateModelStatements(s *schema.Schema) ([]byte, error) {
	packageLine, err := GeneratePackage(s)
	if err != nil {
		return nil, err
	}

	createStatements, err := GenerateCreate(s)
	if err != nil {
		return nil, err
	}

	updateStatements, err := GenerateUpdate(s)
	if err != nil {
		return nil, err
	}

	deleteStatements, err := GenerateDelete(s)
	if err != nil {
		return nil, err
	}

	selectStatements, err := GenerateSelect(s)
	if err != nil {
		return nil, err
	}

	tableName, err := GenerateTableName(s)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(packageLine)
	buf.Write(createStatements)
	buf.Write(updateStatements)
	buf.Write(deleteStatements)
	buf.Write(selectStatements)
	buf.Write(tableName)

	return writers.Clear(buf)
}

// GeneratePackage generates the imports according to the schema.
// TODO remove this function
func GeneratePackage(s *schema.Schema) ([]byte, error) {
	temp := template.New("package.tmpl")
	_, err := temp.Parse(PackageTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// name := strings.ToLower(strings.Split(s.Title, " ")[0])
	name := "models"
	err = temp.ExecuteTemplate(&buf, "package.tmpl", name)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// PackageTemplate holds the template for the packages of the models
var PackageTemplate = `// Generated struct for {{.}}.
package {{.}}
`
