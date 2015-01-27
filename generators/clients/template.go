package clients

// ClientsTemplate holds the template for the clients packages
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

func (a *{{ToUpperFirst .Name}}) One(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.One", req, res)
}

func (a *{{ToUpperFirst .Name}}) Create(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Create", req, res)
}

func (a *{{ToUpperFirst .Name}}) Update(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Update", req, res)
}

func (a *{{ToUpperFirst .Name}}) Delete(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Delete", req, res)
}

func (a *{{ToUpperFirst .Name}}) Some(ctx context.Context, req *models.{{ToUpperFirst .Name}}, res *[]*models.{{ToUpperFirst .Name}}) error {
    return m.client.Call(ctx, "{{ToUpperFirst .Name}}.Some", req, res)
}`
