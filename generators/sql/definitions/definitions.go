package definitions

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

const generatorName = "sql-definition"

type Generator struct {
	DatabaseName  string
	SchemaName    string
	TableName     string
	RoleName      string
	FieldNameCase string `default:"snake"`
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Name() string {
	return generatorName
}

// Generate generates the basic CRUD statements for the models
func (g *Generator) Generate(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	outputs := make([]common.Output, 0)

	if s.Title == "" {
		return outputs, errors.New("Title should be set")
	}

	moduleName := context.FieldNameFunc(s.Title)

	settings := g.generateSettings(moduleName, s)

	for _, name := range schema.SortedKeys(s.Definitions) {
		def := s.Definitions[name]

		// schema should have our generator
		if !def.Generators.Has(g.Name()) {
			continue
		}

		settingsDef := g.setDefaultSettings(settings, def)
		settingsDef.Set("tableName", context.FieldNameFunc(def.Title))

		//
		// generate roles
		//
		role, err := DefineRole(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content:     role,
			Path:        fmt.Sprintf("%s/001-%s_roles.sql", context.Config.Target, settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate database
		//
		db, err := DefineDatabase(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content:     db,
			Path:        fmt.Sprintf("%s/002-%s_database.sql", context.Config.Target, settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate extenstions
		//
		extenstions, err := DefineExtensions(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: extenstions,
			Path: fmt.Sprintf(
				"%s/003-%s_extensions.sql",
				context.Config.Target,
				settingsDef.Get("databaseName").(string)),
			DoNotFormat: true,
		})

		//
		// generate schema
		//
		sc, err := DefineSchema(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: sc,
			Path: fmt.Sprintf(
				"%s/%s/004-schema.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate sequences
		//
		sequence, err := DefineSequence(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: sequence,
			Path: fmt.Sprintf(
				"%s/%s/005-%s-sequence.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
				settingsDef.Get("tableName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate types
		//
		types, err := DefineTypes(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: types,
			Path: fmt.Sprintf(
				"%s/%s/006-%s-types.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
				settingsDef.Get("tableName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate tables
		//
		table, err := DefineTable(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: table,
			Path: fmt.Sprintf(
				"%s/%s/007-%s-table.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
				settingsDef.Get("tableName").(string),
			),
			DoNotFormat: true,
		})

		//
		// generate constraints
		//
		constraints, err := DefineConstraints(context, settingsDef, def)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, common.Output{
			Content: constraints,
			Path: fmt.Sprintf(
				"%s/%s/007-%s-constraints.sql",
				context.Config.Target,
				settingsDef.Get("schemaName").(string),
				settingsDef.Get("tableName").(string),
			),
			DoNotFormat: true,
		})
	}

	return outputs, nil
}

func GenerateDefinitions(context *common.Context, settings schema.Generator, s *schema.Schema) ([]byte, error) {
	context.TemplateFuncs["DefineSQLTable"] = DefineTable
	context.TemplateFuncs["DefineSQLSchema"] = DefineSchema
	context.TemplateFuncs["DefineSQLExtensions"] = DefineExtensions
	context.TemplateFuncs["DefineSQLTypes"] = DefineTypes
	context.TemplateFuncs["DefineSQLSequnce"] = DefineSequence

	temp := template.New("create_statement.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(CreateStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Context  *common.Context
		Schema   *schema.Schema
		Settings schema.Generator
	}{
		Context:  context,
		Schema:   s,
		Settings: settings,
	}

	if err := temp.ExecuteTemplate(&buf, "create_statement.tmpl", data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CreateStatementTemplate holds the template for the create sql statement generator
var CreateStatementTemplate = `{{DefineSQLSchema .Context .Settings .Schema}}

{{DefineSQLSequnce .Context .Settings .Schema}}

{{DefineSQLExtensions .Context .Settings .Schema}}

{{DefineSQLTypes .Context .Settings .Schema}}

{{DefineSQLTable .Context .Settings .Schema}}
`

func GetFieldNameFunc(name string) func(string) string {
	switch name {
	case "lower":
		return strings.ToLower
	default:
		return stringext.ToFieldName
	}
}

func (g *Generator) generateSettings(moduleName string, s *schema.Schema) schema.Generator {
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

func (g *Generator) setDefaultSettings(defaultSettings schema.Generator, s *schema.Schema) schema.Generator {
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
