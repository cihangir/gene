package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineSequence creates definition for sequences
func DefineSequence(settings schema.Generator, s *schema.Schema) ([]byte, error) {
	temp := template.New("create_sequences.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(SequenceTemplate); err != nil {
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
	if err := temp.ExecuteTemplate(&buf, "create_sequences.tmpl", data); err != nil {
		return nil, err
	}

	return clean(buf.Bytes()), nil
}

// SequenceTemplate holds the template for sequences
var SequenceTemplate = `
-- ----------------------------
--  Sequence structure for {{.SchemaName}}.{{.TableName}}_id
-- ----------------------------
DROP SEQUENCE IF EXISTS "{{.SchemaName}}"."{{.TableName}}_id_seq" CASCADE;
CREATE SEQUENCE "{{.SchemaName}}"."{{.TableName}}_id_seq" INCREMENT 1 START 1 MAXVALUE 9223372036854775807 MINVALUE 1 CACHE 1;
GRANT USAGE ON SEQUENCE "{{.SchemaName}}"."{{.TableName}}_id_seq" TO "{{.RoleName}}";
`
