package tests

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {
	if err := GenerateTestFuncs(rootPath); err != nil {
		return err
	}

	if err := GenerateMainTestFileForModule(rootPath, s); err != nil {
		return err
	}

	for _, def := range s.Definitions {
		err := GenerateTests(rootPath, s.Title, def.Title)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateMainTestFileForModule(
	rootPath string,
	s *schema.Schema,
) error {
	// TODO check if file is there, no need to continue

	temp := template.New("mainTestFile.tmpl")
	temp.Funcs(template.FuncMap{
		"ToLowerFirst": stringext.ToLowerFirst,
		"ToLower":      strings.ToLower,
		"ToUpperFirst": stringext.ToUpperFirst,
	})

	_, err := temp.Parse(MainTestsTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "mainTestFile.tmpl", s)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/tests/common_test.go",
		rootPath,
		strings.ToLower(s.Title),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}

func GenerateTestFuncs(rootPath string) error {
	// TODO check if file is there, no need to continue
	temp := template.New("testFuncs.tmpl")
	_, err := temp.Parse(TestFuncs)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "testFuncs.tmpl", nil)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%stests/testfuncs.go", rootPath)

	return writers.WriteFormattedFile(path, buf.Bytes())
}

func GenerateTests(rootPath string, moduleName string, name string) error {
	temp := template.New("tests.tmpl")
	temp.Funcs(template.FuncMap{
		"ToLowerFirst": stringext.ToLowerFirst,
		"ToLower":      strings.ToLower,
		"ToUpperFirst": stringext.ToUpperFirst,
	})

	_, err := temp.Parse(TestsTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	data := struct {
		Name       string
		ModuleName string
	}{
		Name:       name,
		ModuleName: moduleName,
	}

	err = temp.ExecuteTemplate(&buf, "tests.tmpl", data)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sworkers/%s/tests/%s_test.go",
		rootPath,
		moduleName,
		stringext.ToLowerFirst(name),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
