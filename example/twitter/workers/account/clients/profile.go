package accountclient

import (
	"github.com/cihangir/gene/example/twitter/models"
	"github.com/youtube/vitess/go/rpcplus"
	"golang.org/x/net/context"
)

// New creates a new local Profile rpc client
func NewProfile(client *rpcplus.Client) *Profile {
	return &Profile{
		client: client,
	}
}

// Profile is for holding the api functions
type Profile struct {
	client *rpcplus.Client
}

func (p *Profile) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return p.client.Call(ctx, "Profile.Create", req, res)
}

func (p *Profile) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return p.client.Call(ctx, "Profile.Delete", req, res)
}

func (p *Profile) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return p.client.Call(ctx, "Profile.One", req, res)
}

func (p *Profile) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return p.client.Call(ctx, "Profile.Update", req, res)
}
