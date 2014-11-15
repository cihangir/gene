package modules

import (
	"fmt"
	"strings"
	"text/template"

	"bytes"

	"bitbucket.org/cihangirsavas/gene/schema"
	"bitbucket.org/cihangirsavas/gene/stringext"
	"bitbucket.org/cihangirsavas/gene/writers"
)

func (m *Module) GenerateMainFile(rootPath string) error {

	moduleName := stringext.ToLowerFirst(
		m.schema.Title,
	)

	mainFilePath := fmt.Sprintf(
		"%s%s/main.go",
		rootPath,
		fmt.Sprintf(moduleFolderStucture[0], moduleName),
	)

	f, err := generateMainFile(m.schema)
	if err != nil {
		return err
	}

	return writers.WriteFormattedFile(mainFilePath, f)
}

func generateMainFile(s *schema.Schema) ([]byte, error) {
	const templateName = "mainfile.tmpl"
	temp := template.New(templateName)
	temp.Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	})
	_, err := temp.Parse(MainFileTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, templateName, s)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var MainFileTemplate string = `
package main

var (
	Name    = "{{.Title}}"
	VERSION string
)

func main() {
	err := server.New({{ToLower .Title}}api.New()).Listen()
	if err != nil {
		fmt.Println(err.Error)
	}
}
`
