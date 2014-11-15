package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"bitbucket.org/cihangirsavas/gene/stringext"
	"bitbucket.org/cihangirsavas/gene/writers"
)

func Generate(rootPath string, name string) error {
	temp := template.New("handlers.tmpl")
	temp.Funcs(template.FuncMap{
		"ToLowerFirst": stringext.ToLowerFirst,
	})

	_, err := temp.Parse(HandlerTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = temp.ExecuteTemplate(&buf, "handlers.tmpl", name)
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"%sgene/modules/%s/api/%s.go",
		rootPath,
		stringext.ToLowerFirst(name),
		stringext.ToLowerFirst(name),
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
