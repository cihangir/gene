package twitterapi

import (
	"github.com/cihangir/gene/db"
	"github.com/cihangir/gene/example/twitter/models"
	"golang.org/x/net/context"
)

// New creates a new local Tweet handler
func NewTweet() *Tweet { return &Tweet{} }

// Tweet is for holding the api functions
type Tweet struct{}

func (t *Tweet) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Create(models.NewTweet(), req, res)
}

func (t *Tweet) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Delete(models.NewTweet(), req, res)
}

func (t *Tweet) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).One(models.NewTweet(), req, res)
}
