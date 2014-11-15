package handlers

var APITemplate = `package {{ToLowerFirst .}}api

// New creates a new local api handler
func New() gene.Initer { return &api{} }

// api is for holding the interface functions, nothing more
type api struct{}

func (api) Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux {

    // Updates {{ToLowerFirst .}} by it's ID
    mux.Handle("POST", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler: {{ToLowerFirst .}}handlers.Update,
            Name: "{{ToLowerFirst .}}-update",
        },
    ))

    // Deletes {{ToLowerFirst .}} by it's ID
    mux.Handle("DELETE", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}handlers.Delete,
            Name:           "{{ToLowerFirst .}}-delete",
        },
    ))

    // Creates a new {{ToLowerFirst .}} and returns it
    mux.Handle("POST", "/{{ToLowerFirst .}}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}handlers.Create,
            Name:           "{{ToLowerFirst .}}-create",
        },
    ))

    // Get an existing {{ToLowerFirst .}}
    mux.Handle("GET", "/{{ToLowerFirst .}}/{id}", handler.Wrapper(
        handler.Router{
            Handler:        {{ToLowerFirst .}}handlers.Get,
            Name:           "{{ToLowerFirst .}}-get",
        },
    ))

    return mux
}`

var HandlersTemplate = `package {{ToLowerFirst .}}handlers

func Update(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Delete(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Create(u *url.URL, h http.Header, m *models.{{.}}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}

func Get(u *url.URL, h http.Header, _ interface{}, c *models.Context) (int, http.Header, interface{}, error) {
    return 200, nil, nil, nil
}
`
