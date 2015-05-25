package accountclient

import (
	"github.com/fatih/invoker/models"
	"github.com/youtube/vitess/go/rpcplus"
	"golang.org/x/net/context"
)

// New creates a new local Account rpc client
func NewAccount(client *rpcplus.Client) *Account {
	return &Account{
		client: client,
	}
}

// Account is for holding the api functions
type Account struct {
	client *rpcplus.Client
}

func (a *Account) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.Create", req, res)
}

func (a *Account) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.Delete", req, res)
}

func (a *Account) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.One", req, res)
}

func (a *Account) Some(ctx context.Context, req *models.Account, res *[]*models.Account) error {
	return m.client.Call(ctx, "Account.Some", req, res)
}

func (a *Account) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.Update", req, res)
}
