// Package models creates the models for the modules
package models

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihangir/gene/generators/constants"
	"github.com/cihangir/gene/generators/constructors"
	"github.com/cihangir/gene/generators/sql/statements"
	"github.com/cihangir/gene/generators/validators"
	"github.com/cihangir/gene/writers"
	"github.com/cihangir/schema"
	"github.com/cihangir/stringext"
)

// Generate creates the models and write them to the required paths
func Generate(rootPath string, s *schema.Schema) error {

	for _, def := range s.Definitions {
		moduleName := strings.ToLower(def.Title)

		modelFilePath := fmt.Sprintf(
			"%smodels/%s.go",
			rootPath,
			moduleName,
		)

		f, err := GenerateModel(def)
		if err != nil {
			return err
		}

		if err := writers.WriteFormattedFile(modelFilePath, f); err != nil {
			return err
		}
	}
	return nil
}

// GenerateStatements generates the basic CRUD statements for the models
func GenerateStatements(rootPath string, s *schema.Schema) error {

	for _, def := range s.Definitions {
		moduleName := strings.ToLower(def.Title)

		modelFilePath := fmt.Sprintf(
			"%smodels/%s_statements.go",
			rootPath,
			moduleName,
		)

		f, err := GenerateModelStatements(def)
		if err != nil {
			return err
		}

		if err := writers.WriteFormattedFile(modelFilePath, f); err != nil {
			return err
		}
	}
	return nil
}

// GenerateModel generates the model itself
func GenerateModel(s *schema.Schema) ([]byte, error) {
	var buf bytes.Buffer

	packageLine, err := GeneratePackage(s)
	if err != nil {
		return nil, err
	}

	consts, err := constants.Generate(s)
	if err != nil {
		return nil, err
	}

	schema, err := GenerateSchema(s)
	if err != nil {
		return nil, err
	}

	constructor, err := constructors.Generate(s)
	if err != nil {
		return nil, err
	}

	validators, err := validators.Generate(s)
	if err != nil {
		return nil, err
	}

	buf.Write(packageLine)
	buf.Write(consts)
	buf.Write(schema)
	buf.Write(constructor)
	if validators != nil {
		buf.Write(validators)
	}

	return writers.Clear(buf)
}

// GenerateModelStatements generates the CRUD statements for the model struct
func GenerateModelStatements(s *schema.Schema) ([]byte, error) {
	var buf bytes.Buffer

	packageLine, err := GeneratePackage(s)
	if err != nil {
		return nil, err
	}

	createStatements, err := statements.GenerateCreate(s)
	if err != nil {
		return nil, err
	}

	updateStatements, err := statements.GenerateUpdate(s)
	if err != nil {
		return nil, err
	}

	deleteStatements, err := statements.GenerateDelete(s)
	if err != nil {
		return nil, err
	}

	selectStatements, err := statements.GenerateSelect(s)
	if err != nil {
		return nil, err
	}

	tableName, err := statements.GenerateTableName(s)
	if err != nil {
		return nil, err
	}

	buf.Write(packageLine)
	buf.Write(createStatements)
	buf.Write(updateStatements)
	buf.Write(deleteStatements)
	buf.Write(selectStatements)
	buf.Write(tableName)

	return writers.Clear(buf)
}

// GeneratePackage generates the imports according to the schema.
// TODO remove this function
func GeneratePackage(s *schema.Schema) ([]byte, error) {
	temp := template.New("package.tmpl")
	_, err := temp.Parse(PackageTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	// name := strings.ToLower(strings.Split(s.Title, " ")[0])
	name := "models"
	err = temp.ExecuteTemplate(&buf, "package.tmpl", name)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GenerateSchema generates the schema.
func GenerateSchema(s *schema.Schema) ([]byte, error) {

	temp := template.New("schema.tmpl")
	temp.Funcs(schema.Helpers)

	_, err := temp.Parse(StructTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	context := struct {
		Name       string
		Definition *schema.Schema
	}{
		Name:       s.Title,
		Definition: s,
	}

	err = temp.ExecuteTemplate(&buf, "schema.tmpl", context)
	if err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}

// GenerateFunctions generates the functions according to the schema.
func GenerateFunctions(s *schema.Schema) ([]byte, error) {

	temp := template.New("functions.tmpl")
	temp.Funcs(template.FuncMap{
		"Pointerize": stringext.Pointerize,
	})

	_, err := temp.Parse(FunctionsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	context := struct {
		Name  string
		Funcs []string
	}{
		Name: s.Title,
		Funcs: []string{
			"Create", "Update", "Delete", "ById",
			"Some", "One",
		},
	}

	if err := temp.ExecuteTemplate(&buf, "functions.tmpl", context); err != nil {
		return nil, err
	}

	return writers.Clear(buf)
}
