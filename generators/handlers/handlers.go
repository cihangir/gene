package handlers

import (
	"bytes"
	"fmt"
	"text/template"

	"bitbucket.org/cihangirsavas/gene/stringext"
	"bitbucket.org/cihangirsavas/gene/writers"
)

func Generate(name string) error {
	temp := template.New("handlers.tmpl")
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
		"./gene/modules/",
		stringext.ToLowerFirst(name),
		"/api/",
		stringext.ToLowerFirst(name),
		".go",
	)

	return writers.WriteFormattedFile(path, buf.Bytes())
}
