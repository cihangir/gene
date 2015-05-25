package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineTypes creates definition for types
func DefineTypes(settings schema.Generator, s *schema.Schema) ([]byte, error) {
	temp := template.New("create_types.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(TypeTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema     *schema.Schema
		SchemaName string // postgres schema name
		RoleName   string // postgres role name
		TableName  string // postgres table name
	}{
		Schema:     s,
		SchemaName: settings.Get("schemaName").(string),
		RoleName:   settings.Get("roleName").(string),
		TableName:  settings.Get("tableName").(string),
	}
	if err := temp.ExecuteTemplate(&buf, "create_types.tmpl", data); err != nil {
		return nil, err
	}

	return clean(buf.Bytes()), nil
}

// TypeTemplate holds the template for types
var TypeTemplate = `
{{$schemaName := .SchemaName}}
{{$tableName := .TableName}}
{{$roleName := .RoleName}}

{{range $key, $value := .Schema.Properties}}
{{if len $value.Enum}}
-- ----------------------------
--  Types structure for {{$schemaName}}.{{$tableName}}.{{ToFieldName $value.Title}}
-- ----------------------------
DROP TYPE IF EXISTS "{{$schemaName}}"."{{$tableName}}_{{ToFieldName $value.Title}}_enum" CASCADE;
CREATE TYPE "{{$schemaName}}"."{{$tableName}}_{{ToFieldName $value.Title}}_enum" AS ENUM (
  '{{Join $value.Enum "',\n  '"}}'
);
ALTER TYPE "{{$schemaName}}"."{{$tableName}}_{{ToFieldName $value.Title}}_enum" OWNER TO "{{$roleName}}";
{{end}}
{{end}}
`
