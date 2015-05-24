package definitions

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

type generator struct{}

func New() *generator {
	return &generator{}
}

var PathForStatements = "%smodels/%s_statements.go.sql"

func (g *generator) Name() string {
	return "sql-definition"
}

// Generate generates the basic CRUD statements for the models
func (g *generator) Generate(context *config.Context, schema *schema.Schema) ([]common.Output, error) {
	moduleName := context.ModuleNameFunc(schema.Title)
	outputs := make([]common.Output, 0)

	for _, def := range schema.Definitions {

		// schema should have our generator
		if !def.Generators.Has(g.Name()) {
			continue
		}

		settings, _ := def.Generators.Get(g.Name())
		settings.SetNX("schemaName", stringext.ToFieldName(moduleName))
		settings.SetNX("tableName", stringext.ToFieldName(def.Title))
		settings.SetNX("roleName", stringext.ToFieldName(moduleName))

		// convert []interface to []string
		grants := settings.GetWithDefault("grants", []string{"ALL"})
		grantsI, ok := grants.([]interface{})
		grantsS := make([]string, 0)

		if ok {
			for _, t := range grantsI {
				grantsS = append(grantsS, t.(string))
			}
		} else {
			grantsS = grants.([]string)
		}

		settings.Set("grants", grantsS)

		f, err := GenerateDefinitions(settings, def)
		if err != nil {
			return outputs, err
		}

		path := fmt.Sprintf(PathForStatements, context.Config.Target, moduleName)

		outputs = append(outputs, common.Output{
			Content:     f,
			Path:        path,
			DoNotFormat: true,
		})
	}

	return outputs, nil
}

func GenerateDefinitions(settings schema.Generator, s *schema.Schema) ([]byte, error) {
	common.TemplateFuncs["DefineSQLTable"] = DefineTable
	common.TemplateFuncs["DefineSQLSchema"] = DefineSchema
	common.TemplateFuncs["DefineSQLExtensions"] = DefineExtensions
	common.TemplateFuncs["DefineSQLTypes"] = DefineTypes
	common.TemplateFuncs["DefineSQLSequnce"] = DefineSequence

	temp := template.New("create_statement.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(CreateStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema   *schema.Schema
		Settings schema.Generator
	}{
		Schema:   s,
		Settings: settings,
	}
	if err := temp.ExecuteTemplate(&buf, "create_statement.tmpl", data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func clean(b []byte) []byte {
	b = writers.NewLinesRegex.ReplaceAll(b, []byte(""))

	// convert tabs to 4 spaces
	b = bytes.Replace(b, []byte("\t"), []byte("    "), -1)

	// clean extra spaces
	b = bytes.Replace(b, []byte("  ,"), []byte(","), -1)
	b = bytes.Replace(b, []byte(" ,"), []byte(","), -1)

	// replace last trailing comma
	b = bytes.Replace(b, []byte(",\n)"), []byte("\n)"), -1)

	return b
}

// CreateStatementTemplate holds the template for the create sql statement generator
var CreateStatementTemplate = `{{DefineSQLSchema .Settings .Schema}}

{{DefineSQLSequnce .Settings .Schema}}

{{DefineSQLExtensions .Settings .Schema}}

{{DefineSQLTypes .Settings .Schema}}

{{DefineSQLTable .Settings .Schema}}
`
