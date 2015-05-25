package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineRoles creates definition for types
func DefineRoles(settings schema.Generator, s *schema.Schema) (res string) {
	temp := template.New("create_role.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(RoleTemplate); err != nil {
		panic(err)
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
		panic(err)
	}

	return string(clean(buf.Bytes()))
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
