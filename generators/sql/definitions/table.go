package definitions

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

// DefineSequence creates definition for sequences
func DefineTable(settings schema.Generator, s *schema.Schema) (res string) {
	common.TemplateFuncs["GenerateSQLField"] = GenerateSQLField

	temp := template.New("create_table.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(TableTemplate); err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	data := struct {
		Schema     *schema.Schema
		SchemaName string // postgres schema name
		RoleName   string // postgres role name
		TableName  string // postgres table name
		Settings   schema.Generator
	}{
		Schema:     s,
		SchemaName: settings.Get("schemaName").(string),
		RoleName:   settings.Get("roleName").(string),
		TableName:  settings.Get("tableName").(string),
		Settings:   settings,
	}
	if err := temp.ExecuteTemplate(&buf, "create_table.tmpl", data); err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}

// TableTemplate holds the template for sequences
var TableTemplate = `
{{$schema := .Schema}}
{{$settings := .Settings}}
-- ----------------------------
--  Table structure for {{.SchemaName}}.{{.TableName}}
-- ----------------------------
DROP TABLE IF EXISTS "{{.SchemaName}}"."{{.TableName}}";
CREATE TABLE "{{.SchemaName}}"."{{.TableName}}" (
{{range $key, $value := .Schema.Properties}}
    {{GenerateSQLField $settings $value $key}}
{{end}}
) WITH (OIDS = FALSE);-- end schema creation

GRANT {{Join $settings.grants ", "}} ON "{{.SchemaName}}"."{{.TableName}}" TO "{{$settings.roleName}}";
`

// DefineTable creates a definition line for a given coloumn
func GenerateSQLField(settings schema.Generator, s *schema.Schema, propertyName string) (res string) {
	schemaName := settings.Get("schemaName").(string)
	tableName := settings.Get("tableName").(string)

	property := s

	fieldName := stringext.ToFieldName(propertyName) // transpiled version of property
	if property.Title != "" {
		fieldName = stringext.ToFieldName(property.Title)
	}

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
