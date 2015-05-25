package definitions

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/config"
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/generators/folders"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

const generatorName = "sql-definition"

type generator struct{}

func New() *generator {
	return &generator{}
}

var PathForStatements = "%smodels/%s_statements.go.sql"

func (g *generator) Name() string {
	return generatorName
}

func (g *generator) generateSettings(moduleName string, s *schema.Schema) schema.Generator {
	settings, ok := s.Generators.Get(g.Name())
	if !ok {
		settings = schema.Generator{}
	}
	settings.SetNX("databaseName", stringext.ToFieldName(moduleName))
	settings.SetNX("schemaName", stringext.ToFieldName(moduleName))
	settings.SetNX("tableName", stringext.ToFieldName(s.Title))
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

	return settings
}

func (g *generator) setDefaultSettings(defaultSettings schema.Generator, s *schema.Schema) schema.Generator {
	settings, _ := s.Generators.Get(g.Name())

	settings.SetNX("databaseName", defaultSettings.Get("databaseName").(string))
	settings.SetNX("schemaName", defaultSettings.Get("schemaName").(string))
	settings.SetNX("tableName", defaultSettings.Get("tableName").(string))
	settings.SetNX("roleName", defaultSettings.Get("roleName").(string))

	// convert []interface to []string
	grants := settings.GetWithDefault("grants", defaultSettings.Get("grants").([]string))
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

	return settings
}

// Generate generates the basic CRUD statements for the models
func (g *generator) Generate(context *config.Context, schema *schema.Schema) ([]common.Output, error) {
	outputs := make([]common.Output, 0)

	if schema.Title == "" {
		return outputs, errors.New("Title should be set")
	}

	moduleName := context.ModuleNameFunc(schema.Title)

	settings := g.generateSettings(moduleName, schema)

	for _, def := range schema.Definitions {

		// schema should have our generator
		if !def.Generators.Has(g.Name()) {
			continue
		}

		settingsDef := g.setDefaultSettings(settings, def)
		settingsDef.Set("tableName", stringext.ToFieldName(def.Title))

		//
		// generate database
		//
		db, err := DefineDatabase(settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content:     db,
			Path:        fmt.Sprintf("%sdb/001-%s_database.sql", context.Config.Target, settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate roles
		//
		role, err := DefineRole(settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content:     role,
			Path:        fmt.Sprintf("%sdb/002-%s_roles.sql", context.Config.Target, settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate extenstions
		//
		extenstions, err := DefineExtensions(settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content:     extenstions,
			Path:        fmt.Sprintf("%sdb/003-%s_extensions.sql", context.Config.Target, settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate sequences
		//
		sequence, err := DefineSequence(settingsDef, def)
		if err != nil {
			return nil, err
		}

		// create the module folder structure
		if err := folders.EnsureFolders(
			context.Config.Target, // root folder
			[]string{fmt.Sprintf(
				"db/%s", settingsDef.Get("schemaName").(string),
			)},
		); err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: sequence,
			Path: fmt.Sprintf(
				"%sdb/%s/004-sequence.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate types
		//
		types, err := DefineTypes(settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: types,
			Path: fmt.Sprintf(
				"%sdb/%s/005-types.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate tables
		//
		table, err := DefineTable(settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: table,
			Path: fmt.Sprintf(
				"%sdb/%s/006-table.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
			),
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
