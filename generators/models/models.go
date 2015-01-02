package models

import (
	"bytes"
	"fmt"
	"text/template"

	"bitbucket.org/cihangirsavas/gene/generators/validators"
	"bitbucket.org/cihangirsavas/gene/schema"
	"bitbucket.org/cihangirsavas/gene/stringext"
	"bitbucket.org/cihangirsavas/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {
	moduleName := stringext.ToLowerFirst(
		s.Title,
	)

	modelFilePath := fmt.Sprintf(
		"%sgene/models/%s.go",
		rootPath,
		moduleName,
	)

	f, err := GenerateModel(s)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(modelFilePath, f)
}

func GenerateModel(s *schema.Schema) ([]byte, error) {
	var buf bytes.Buffer

	packageLine, err := GeneratePackage(s)
	if err != nil {
		return nil, err
	}

	schema, err := GenerateSchema(s)
	if err != nil {
		return nil, err
	}

	// validators, err := GenerateValidators(s)
	// if err != nil {
	// 	return nil, err
	// }

	// funcs, err := GenerateFunctions(s)
	// if err != nil {
	// 	return nil, err
	// }

	buf.WriteString(string(packageLine))
	buf.WriteString(string(schema))
	// buf.WriteString(string(validators))
	// buf.WriteString(string(funcs))

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
func GenerateValidators(s *schema.Schema) ([]byte, error) {
	temp := template.New("validators.tmpl")
	temp.Funcs(template.FuncMap{
		"GenerateValidator": validators.GenerateValidator,
		"Pointerize":        stringext.Pointerize,
	})

	_, err := temp.Parse(ValidatorsTemplate)
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

	temp.ExecuteTemplate(&buf, "validators.tmpl", context)

	return writers.Clear(buf)
}

// Generate functions according to the schema.
func GenerateFunctions(s *schema.Schema) ([]byte, error) {

	temp := template.New("functions.tmpl")
	temp.Funcs(template.FuncMap{
		"GenerateValidator": validators.GenerateValidator,
		"Pointerize":        stringext.Pointerize,
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
