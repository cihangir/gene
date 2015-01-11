package models

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/cihangir/gene/generators/constants"
	"github.com/cihangir/gene/generators/constructors"
	"github.com/cihangir/gene/generators/sql/statements"
	"github.com/cihangir/gene/generators/validators"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {

	for _, def := range s.Definitions {
		moduleName := stringext.ToLowerFirst(def.Title)

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

	statements, err := statements.Generate(s)
	if err != nil {
		return nil, err
	}

	buf.WriteString(string(packageLine))
	buf.WriteString(string(consts))
	buf.WriteString(string(schema))
	buf.WriteString(string(constructor))
	if validators != nil {
		buf.WriteString(string(validators))
	}

	buf.WriteString(string(statements))

	return writers.Clear(buf)
}

// Generate imports according to the schema.
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

// Generate schema according to the schema.
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

// Generate functions according to the schema.
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
