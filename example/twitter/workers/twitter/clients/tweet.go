package twitterclient

import (
	"github.com/fatih/invoker/models"
	"github.com/youtube/vitess/go/rpcplus"
	"golang.org/x/net/context"
)

// New creates a new local Tweet rpc client
func NewTweet(client *rpcplus.Client) *Tweet {
	return &Tweet{
		client: client,
	}
}

// Tweet is for holding the api functions
type Tweet struct {
	client *rpcplus.Client
}

func (t *Tweet) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return t.client.Call(ctx, "Tweet.Create", req, res)
}

func (t *Tweet) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return t.client.Call(ctx, "Tweet.Delete", req, res)
}

func (t *Tweet) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return t.client.Call(ctx, "Tweet.One", req, res)
}
