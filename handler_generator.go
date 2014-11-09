package gene

import (
	"bytes"
	"fmt"
)

func (m *Module) GenerateHandlers() error {
	temp := templates.New("handlers.tmpl")
	_, err := temp.Parse(HandlerTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = templates.ExecuteTemplate(&buf, "handlers.tmpl", m.schema)
	if err != nil {
		return err
	}

	fmt.Println("buf.String()-->", buf.String())
	return err
}

var HandlerTemplate = `
// Updates {{lowFirst .Title}} by it's ID
mux.Handle("POST", "/{{lowFirst .Title}}/{id}", handler.Wrapper(
    handler.Request{
        Handler:        {{lowFirst .Title}}.Update,
        Name:           "{{lowFirst .Title}}-update",
    },
))

// Deletes {{lowFirst .Title}} by it's ID
mux.Handle("DELETE", "/{{lowFirst .Title}}/{id}", handler.Wrapper(
    handler.Request{
        Handler:        {{lowFirst .Title}}.Delete,
        Name:           "{{lowFirst .Title}}-delete",
    },
))

// Creates a new {{lowFirst .Title}} and returns it
mux.Handle("POST", "/{{lowFirst .Title}}", handler.Wrapper(
    handler.Request{
        Handler:        {{lowFirst .Title}}.Create,
        Name:           "{{lowFirst .Title}}-create",
    },
))
`
