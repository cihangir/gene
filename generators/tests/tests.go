package tests

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/schema"
	"github.com/cihangir/gene/stringext"
	"github.com/cihangir/gene/writers"
)

func Generate(rootPath string, s *schema.Schema) error {
	// Generate test functions
	testFuncs, err := GenerateTestFuncs(s)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%stests/testfuncs.go", rootPath)

	if err := writers.WriteFormattedFile(path, testFuncs); err != nil {
		return err
	}

	// generate module test file
	mainTest, err := GenerateMainTestFileForModule(s)
	if err != nil {
		return err
	}

	path = fmt.Sprintf(
		"%sworkers/%s/tests/common_test.go",
		rootPath,
		strings.ToLower(s.Title),
	)

	if err := writers.WriteFormattedFile(path, mainTest); err != nil {
		return err
	}

	// generate tests for the schema
	for _, def := range s.Definitions {
		testFile, err := GenerateTests(s.Title, def.Title)
		if err != nil {
			return err
		}
		path = fmt.Sprintf(
			"%sworkers/%s/tests/%s_test.go",
			rootPath,
			s.Title,
			stringext.ToLowerFirst(def.Title),
		)

		return writers.WriteFormattedFile(path, testFile)
	}

	return nil
}

func GenerateMainTestFileForModule(s *schema.Schema) ([]byte, error) {
	// TODO check if file is there, no need to continue

	temp := template.New("mainTestFile.tmpl").Funcs(common.TemplateFuncs)

	_, err := temp.Parse(MainTestsTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "mainTestFile.tmpl", s)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func GenerateTestFuncs(s *schema.Schema) ([]byte, error) {
	// TODO check if file is there, no need to continue
	temp := template.New("testFuncs.tmpl")
	_, err := temp.Parse(TestFuncs)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "testFuncs.tmpl", nil)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}

func GenerateTests(moduleName string, name string) ([]byte, error) {
	temp := template.New("tests.tmpl").Funcs(common.TemplateFuncs)
	_, err := temp.Parse(TestsTemplate)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return format.Source(buf.Bytes())
}
