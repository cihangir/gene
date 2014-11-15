package handlers

var HandlerTemplate = `package {{lowFirst .}}api

// New creates a new local api handler
func New() gene.Initer { return &api{} }

// api is for holding the interface functions, nothing more
type api struct{}

func (api) Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux {

    // Updates {{lowFirst .}} by it's ID
    mux.Handle("POST", "/{{lowFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler: {{lowFirst .}}.Update,
            Name: "{{lowFirst .}}-update",
        },
    ))

    // Deletes {{lowFirst .}} by it's ID
    mux.Handle("DELETE", "/{{lowFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler:        {{lowFirst .}}.Delete,
            Name:           "{{lowFirst .}}-delete",
        },
    ))

    // Creates a new {{lowFirst .}} and returns it
    mux.Handle("POST", "/{{lowFirst .}}", handler.Wrapper(
        handler.Router{
            Handler:        {{lowFirst .}}.Create,
            Name:           "{{lowFirst .}}-create",
        },
    ))
}`
