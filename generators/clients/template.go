package clients

var ClientsTemplate = `package {{ToLowerFirst .}}client

import (
    "github.com/youtube/vitess/go/rpcplus"
    "golang.org/x/net/context"
)

// New creates a new {{.}} rpc client
func New(client *rpcplus.Client) *{{.}} {
    return &{{.}}{
        client: client,
    }
}

// {{.}} is for holding the api functions
type {{.}} struct{
    client *rpcplus.Client
}

// generate this for all indexes
// func (m *{{.}}) ById(ctx context.Context, id *int64, res *models.{{.}}) error {
//   return m.client.Call(ctx, "{{.}}.ById", id, res)
// }

// generate this for all indexes
// func (m *{{.}}) ByIds(ctx context.Context, ids *[]int64, res *[]*models.{{.}}) error {
//   return m.client.Call(ctx, "{{.}}.ByIds", id, res)
// }

func (m *{{.}}) One(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return m.client.Call(ctx, "{{.}}.One", req, res)
}

func (m *{{.}}) Create(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return m.client.Call(ctx, "{{.}}.Create", req, res)
}

func (m *{{.}}) Update(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return m.client.Call(ctx, "{{.}}.Update", req, res)
}

func (m *{{.}}) Delete(ctx context.Context, req *models.{{.}}, res *models.{{.}}) error {
    return m.client.Call(ctx, "{{.}}.Delete", req, res)
}

func (m *{{.}}) Some(ctx context.Context, req *request.Options, res *[]*models.{{.}}) error {
    return m.client.Call(ctx, "{{.}}.Some", req, res)
}`
