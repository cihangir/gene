package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineRole creates definition for types
func DefineRole(settings schema.Generator, s *schema.Schema) ([]byte, error) {
	temp := template.New("create_role.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(RoleTemplate); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	data := struct {
		Schema       *schema.Schema
		DatabaseName string // postgres database name
		SchemaName   string // postgres schema name
		RoleName     string // postgres role name
		TableName    string // postgres table name
	}{
		Schema:       s,
		DatabaseName: settings.Get("databaseName").(string),
		SchemaName:   settings.Get("schemaName").(string),
		RoleName:     settings.Get("roleName").(string),
		TableName:    settings.Get("tableName").(string),
	}
	if err := temp.ExecuteTemplate(&buf, "create_role.tmpl", data); err != nil {
		return nil, err
	}

	return clean(buf.Bytes()), nil
}

// RoleTemplate holds the template for types
var RoleTemplate = `
{{$databaseName := .DatabaseName}}
{{$roleName := .RoleName}}

-- Drop role
DROP ROLE IF EXISTS "{{$roleName}}";

-- Create role
CREATE ROLE "{{$roleName}}";
`
