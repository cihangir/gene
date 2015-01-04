package clients

var ClientsTemplate = `package {{ToLower .ModuleName}}client

import (
    "github.com/youtube/vitess/go/rpcplus"
    "golang.org/x/net/context"
)

// New creates a new {{.}} rpc client
func New{{ToUpperFirst .Name}}(client *rpcplus.Client) *{{ToUpperFirst .Name}} {
    return &{{ToUpperFirst .Name}}{
        client: client,
    }
}

// {{ToUpperFirst .Name}} is for holding the api functions
type {{ToUpperFirst .Name}} struct{
    client *rpcplus.Client
}

// generate this for all indexes
// func (m *{{ToUpperFirst .Name}}) ById(ctx context.Context, id *int64, res *models.{{ToUpperFirst .Name}}) error {
//   return m.client.Call(ctx, "{{ToUpperFirst .Name}}.ById", id, res)
// }

// generate this for all indexes
// func (m *{{ToUpperFirst .Name}}) ByIds(ctx context.Context, ids *[]int64, res *[]*models.{{ToUpperFirst .Name}}) error {
//   return m.client.Call(ctx, "{{ToUpperFirst .Name}}.ByIds", id, res)
// }

func (m *{{ToUpperFirst .Name}}) One(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.One", req, res)
}

func (m *{{ToUpperFirst .Name}}) Create(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Create", req, res)
}

func (m *{{ToUpperFirst .Name}}) Update(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Update", req, res)
}

func (m *{{ToUpperFirst .Name}}) Delete(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Delete", req, res)
}

func (m *{{ToUpperFirst .Name}}) Some(ctx context.Context, req *request.Options, res *[]*models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Some", req, res)
}`
