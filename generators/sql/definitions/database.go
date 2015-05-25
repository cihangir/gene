package definitions

import (
	"bytes"
	"text/template"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

// DefineDatabase creates definition for types
func DefineDatabase(settings schema.Generator, s *schema.Schema) (res string) {
	temp := template.New("create_database.tmpl").Funcs(common.TemplateFuncs)
	if _, err := temp.Parse(DatabaseTemplate); err != nil {
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
	if err := temp.ExecuteTemplate(&buf, "create_database.tmpl", data); err != nil {
		panic(err)
	}

	return string(clean(buf.Bytes()))
}

//  DatabaseTemplate holds the template for types
var DatabaseTemplate = `
{{$databaseName := .DatabaseName}}
{{$roleName := .RoleName}}

-- Drop database
DROP DATABASE IF EXISTS "{{$databaseName}}";

-- Create database itself
CREATE DATABASE "{{$databaseName}}" OWNER "{{$roleName}}" ENCODING 'UTF8'  TEMPLATE template0;
`