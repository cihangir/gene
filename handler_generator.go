package gene

import "bytes"

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

	path := "./gene/modules/" + lowFirst(m.schema.Title) + "/api/" + lowFirst(m.schema.Title) + ".go"

	return writeFormattedFile(path, buf.Bytes())
}

var HandlerTemplate = `package {{lowFirst .Title}}api

// New creates a new local api handler
func New() gene.Initer { return &api{} }

// api is for holding the interface functions, nothing more
type api struct{}

func (api) Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux {

    // Updates {{lowFirst .Title}} by it's ID
    mux.Handle("POST", "/{{lowFirst .Title}}/{id}", handler.Wrapper(
        handler.Request{
            Handler: {{lowFirst .Title}}.Update,
            Name: "{{lowFirst .Title}}-update",
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
}`
