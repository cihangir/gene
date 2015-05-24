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
	}{
		Schema:     s,
		SchemaName: settings.Get("schemaName").(string),
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
CREATE SCHEMA IF NOT EXISTS "{{.SchemaName}}"
`
