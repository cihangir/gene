package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineSchema creates definition for schema
func DefineSchema(settings schema.Generator, s *schema.Schema) (res string) {
	temp := template.New("create_schema.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(SchemaTemplate); err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	data := struct {
		Schema     *schema.Schema
		SchemaName string // postgres schema name
		RoleName   string
	}{
		Schema:     s,
		SchemaName: settings.Get("schemaName").(string),
		RoleName:   settings.Get("roleName").(string),
	}
	if err := temp.ExecuteTemplate(&buf, "create_schema.tmpl", data); err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}

//  SchemaTemplate holds the template for sequences
var SchemaTemplate = `
-- ----------------------------
--  Schema structure for {{.SchemaName}}
-- ----------------------------
-- create schema
CREATE SCHEMA IF NOT EXISTS "{{.SchemaName}}";

-- give usage permission
GRANT usage ON SCHEMA "{{.SchemaName}}" to "{{.RoleName}}";

-- add new schema to search path -just for convenience
-- SELECT set_config('search_path', current_setting('search_path') || ',{{.SchemaName}}', false);
`
