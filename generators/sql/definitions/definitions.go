package definitions

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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
		settings.SetNX("roleName", stringext.ToFieldName(moduleName))

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
	common.TemplateFuncs["DefineSQLField"] = Define

	temp := template.New("create_statement.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(CreateStatementTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema     *schema.Schema
		SchemaName string // postgres schema name
		RoleName   string // postgres role name
	}{
		Schema:     s,
		SchemaName: settings.Get("schemaName"),
		RoleName:   settings.Get("roleName"),
	}
	if err := temp.ExecuteTemplate(&buf, "create_statement.tmpl", data); err != nil {
		return nil, err
	}

	return clean(buf.Bytes()), nil
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
var CreateStatementTemplate = `
{{$schemaName := ToFieldName .SchemaName}}
{{$title := ToFieldName .Schema.Title}}
{{$schema := .Schema}}


-- ----------------------------
--  Sequence structure for {{$schemaName}}.{{$title}}_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "{{$schemaName}}"."{{$title}}_id_seq";
CREATE SEQUENCE "{{$schemaName}}"."{{$title}}_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "{{$schemaName}}"."{{$title}}_id_seq" TO "{{.RoleName}}";

-- ----------------------------
--  Table structure for {{$schemaName}}.{{$title}}
-- ----------------------------
DROP TABLE IF EXISTS "{{$schemaName}}"."{{$title}}";
CREATE TABLE "{{$schemaName}}"."{{$title}}" (
{{range $key, $value := .Schema.Properties}}
    {{DefineSQLField $schemaName $title $key $schema}}
{{end}}
) WITH (OIDS = FALSE);-- end schema creation
`

// Define creates a definition line for a given coloumn
func Define(schemaName string, tableName string, propertyName string, s *schema.Schema) (res string) {

	property := s.Properties[propertyName]

	fieldName := stringext.ToFieldName(propertyName) // transpiled version of property

	fieldType := "" // will hold the type for coloumn

	switch strings.ToLower(property.Type.(string)) {
	case "boolean":
		fieldType = "BOOLEAN"
	case "string":
		switch property.Format {
		case "date-time":
			fieldType = "TIMESTAMP (6) WITH TIME ZONE"
		case "UUID":
			fieldType = "UUID"
		default:
			typeName := "TEXT"
			if property.MaxLength > 0 {
				// if schema defines a max length, no need to use text
				typeName = fmt.Sprintf("VARCHAR (%d)", property.MaxLength)
			}

			fieldType = fmt.Sprintf("%s COLLATE \"default\"", typeName)
		}
	case "number":
		fieldType = "NUMERIC"

		switch property.Format {
		case "int64", "uint64":
			fieldType = "BIGINT"
		case "integer", "int", "int32", "uint", "uint32":
			fieldType = "INTEGER"
		case "int8", "uint8", "int16", "uint16":
			fieldType = "SMALLINT"
		case "float32", "float64":
			fieldType = "NUMERIC"
		}
	case "any":
		panic("should specify type")
	case "array":
		panic("should specify type")
	case "object", "config":
		// TODO implement embedded struct table creation
		res = ""
	case "null":
		res = ""
	case "error":
		res = ""
	case "custom":
		res = ""
	default:
		panic("unknown field")
	}

	// override if it is an enum field
	if len(property.Enum) > 0 {
		fieldType = fmt.Sprintf(
			"\"%s_%s_enum\"",
			stringext.ToFieldName(tableName),
			stringext.ToFieldName(propertyName),
		)
	}

	res = fmt.Sprintf(
		"%q %s %s %s %s,",
		fieldName, // first name comes
		fieldType, // then type of the coloumn
		generateDefaultValue(schemaName, fieldName, tableName, property), // generate default value if exists
		generateNotNull(s, propertyName),                                 // generate not null statement if requiired
		generateCheckStatements(tableName, fieldName, property),          // generate validators
	)

	return res
}

// generateDefaultValue generates `default` string for given coloumn
func generateDefaultValue(schemaName string, propertyName, tableName string, s *schema.Schema) string {
	if propertyName == "id" {
		return fmt.Sprintf("DEFAULT nextval('%s.%s_id_seq' :: regclass) ", schemaName, tableName)
	}

	if s.Default == nil {
		return ""
	}

	if len(s.Enum) > 0 {
		// enums should be a valud enum string
		if !common.IsIn(s.Default.(string), s.Enum...) {
			fmt.Printf("%s not a valid enum", s.Default)
			os.Exit(1)
		}

		return fmt.Sprintf("DEFAULT %q", s.Default)
	}

	def := ""
	switch s.Default.(type) {
	case float64, float32, int16, int32, int, int64, uint16, uint32, uint, uint64, bool:
		def = fmt.Sprintf("%v", s.Default)
	default:
		def = fmt.Sprintf("%v", s.Default)
		if strings.HasSuffix(def, "()") {
			return fmt.Sprintf("DEFAULT %s", def)
		} else {
			def = fmt.Sprintf("'%v'", s.Default)
		}
	}

	return fmt.Sprintf("DEFAULT %s", strings.ToUpper(def))
}

// generateNotNull if field is int required values, set NOT NULL
func generateNotNull(s *schema.Schema, name string) string {
	for _, n := range s.Required {
		if name == n {
			return "NOT NULL"
		}
	}

	return ""
}

// generateCheckStatements generates validators
func generateCheckStatements(tableName, fieldName string, property *schema.Schema) string {
	chekcs := ""
	switch strings.ToLower(property.Type.(string)) {
	case "string":
		if property.MinLength > 0 {
			chekcs += fmt.Sprintf(
				"\n\t\tCONSTRAINT \"check_%s_%s_min_length_%d\" CHECK (char_length(%q) > %d )",
				tableName,
				fieldName,
				property.MinLength,
				fieldName,
				property.MinLength,
			)
		}

		if property.Pattern != "" {
			chekcs += fmt.Sprintf(
				"\n\t\tCONSTRAINT \"check_%s_%s_pattern\" CHECK (%q ~ '%s')",
				tableName,
				fieldName,
				fieldName,
				property.Pattern,
			)
		}
		// no need to check for max length, we already create coloumn with max length
	case "number":
		if property.MultipleOf > 0 {
			chekcs += fmt.Sprintf(
				"\n\t\tCONSTRAINT \"check_%s_%s_multiple_of_%d\" CHECK ((%q %% %f) = 0)",
				tableName,
				fieldName,
				int64(property.MultipleOf), // do not use dot in check constraint
				fieldName,
				property.MultipleOf,
			)
		}

		if property.Maximum > 0 {
			checker := "<"
			str := "lt"

			if !property.ExclusiveMaximum {
				checker += "="
				str += "e"
			}

			chekcs += fmt.Sprintf(
				"\n\t\tCONSTRAINT \"check_%s_%s_%s_%d\" CHECK (%q %s %f)",
				tableName,
				fieldName,
				str,
				int64(property.Maximum),
				fieldName,
				checker,
				property.Maximum,
			)
		}

		if property.Minimum > 0 {
			checker := ">"
			str := "gt"

			if !property.ExclusiveMinimum {
				checker += "="
				str += "e"
			}

			chekcs += fmt.Sprintf(
				"\n\t\tCONSTRAINT \"check_%s_%s_%s_%d\" CHECK (%q %s %f)",
				tableName,
				fieldName,
				str,
				int64(property.Maximum),
				fieldName,
				checker,
				property.Maximum,
			)
		}
	}

	return chekcs
}
