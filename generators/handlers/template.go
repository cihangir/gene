package handlers

var HandlerTemplate = `package {{ToLowerFirst .}}api

// New creates a new local api handler
func New() gene.Initer { return &api{} }

// api is for holding the interface functions, nothing more
type api struct{}

func (api) Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux {

    // Updates {{ToLowerFirst .}} by it's ID
    mux.Handle("POST", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler: {{ToLowerFirst .}}.Update,
            Name: "{{ToLowerFirst .}}-update",
        },
    ))

    // Deletes {{ToLowerFirst .}} by it's ID
    mux.Handle("DELETE", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}.Delete,
            Name:           "{{ToLowerFirst .}}-delete",
        },
    ))

    // Creates a new {{ToLowerFirst .}} and returns it
    mux.Handle("POST", "/{{ToLowerFirst .}}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}.Create,
            Name:           "{{ToLowerFirst .}}-create",
        },
    ))

    // Get an existing {{ToLowerFirst .}}
    mux.Handle("GET", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}.Get,
            Name:           "{{ToLowerFirst .}}-get",
        },
    ))

    return mux
}`
